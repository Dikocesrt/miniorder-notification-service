package utils_test

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"testing"

	"ordermini-notification-service/internal/domain"
	"ordermini-notification-service/pkg/utils"
)

func runDummySMTPServer(t *testing.T) (host string, port int, closeFunc func()) {
	// Meminta OS mencarikan port acak yang kosong
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("gagal menjalankan dummy tcp listen: %v", err)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				scanner := bufio.NewScanner(c)

				// Sapaan selamat datang awal dari Server SMTP
				fmt.Fprint(c, "220 mock-smtp\r\n")

				for scanner.Scan() {
					line := scanner.Text()

					// Menandakan akhir dari pengiriman blok isi pesan (Body) email.
					if line == "." {
						fmt.Fprint(c, "250 ok\r\n")
						continue
					}

					parts := strings.Fields(line)
					if len(parts) == 0 {
						continue
					}

					cmd := strings.ToUpper(parts[0])
					switch cmd {
					case "EHLO", "HELO":
						fmt.Fprint(c, "250-mock-smtp\r\n250 AUTH PLAIN\r\n")
					case "AUTH":
						fmt.Fprint(c, "235 ok\r\n")
					case "MAIL":
						fmt.Fprint(c, "250 ok\r\n")
					case "RCPT":
						fmt.Fprint(c, "250 ok\r\n")
					case "DATA":
						fmt.Fprint(c, "354 Go ahead\r\n")
					case "QUIT":
						fmt.Fprint(c, "221 ok\r\n")
						return
					}
				}
			}(conn)
		}
	}()

	addr := l.Addr().(*net.TCPAddr)
	return addr.IP.String(), addr.Port, func() { l.Close() }
}

func TestSendSuccessPaymentEmail(t *testing.T) {
	// Menjalankan server SMTP bohongan di Background khusus untuk unit test ini
	host, port, closeSrv := runDummySMTPServer(t)
	defer closeSrv()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.Background()

	t.Run("Berhasil mengirim email - Happy Path", func(t *testing.T) {
		svcProps := domain.ServiceProperties{
			SmtpConfig: domain.SmtpConfig{
				Host:     host,
				Port:     port,
				Username: "test@gmail.com",
				Password: "mock-password-123",
			},
		}

		err := utils.SendSuccessPaymentEmail(ctx, logger, svcProps, "ORDER-123XYZ", "user@example.com")
		if err != nil {
			t.Errorf("Seharusnya tidak ada error, namun mendapatkan: %v", err)
		}
	})

	t.Run("Gagal karena konfigurasi kredensial kosong", func(t *testing.T) {
		svcProps := domain.ServiceProperties{
			SmtpConfig: domain.SmtpConfig{
				Host: "",
				Port: 0,
			},
		}

		err := utils.SendSuccessPaymentEmail(ctx, logger, svcProps, "ORDER-123XYZ", "user@example.com")
		if err == nil {
			t.Errorf("Seharusnya mengembalikan error karena config kosong, tetapi nil")
		}
	})

	t.Run("Gagal melempar ke Mail Server terdekat (Connection Refused)", func(t *testing.T) {
		svcProps := domain.ServiceProperties{
			SmtpConfig: domain.SmtpConfig{
				Host:     host,
				Port:     12345, // Port sembarangan / ditutup
				Username: "test@gmail.com",
				Password: "mock-password",
			},
		}

		err := utils.SendSuccessPaymentEmail(ctx, logger, svcProps, "ORDER-123XYZ", "user@example.com")
		if err == nil {
			t.Errorf("Seharusnya error akibat port yang salah, tetapi nil")
		}
	})
}
