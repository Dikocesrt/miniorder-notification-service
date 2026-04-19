package utils

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
	"ordermini-notification-service/internal/domain"
)

func SendSuccessPaymentEmail(ctx context.Context, logger *slog.Logger, svcProperties domain.ServiceProperties, orderID string, email string) error {
	logger.InfoContext(ctx, "menyiapkan pengiriman email (SMTP Gmail)...",
		slog.String("order_id", orderID),
		slog.String("recipient_email", email),
	)

	smtpHost := svcProperties.SmtpConfig.Host
	smtpPort := svcProperties.SmtpConfig.Port
	smtpAuthEmail := svcProperties.SmtpConfig.Username
	smtpAuthPassword := svcProperties.SmtpConfig.Password

	if smtpHost == "" || smtpPort == 0 || smtpAuthEmail == "" || smtpAuthPassword == "" {
		return fmt.Errorf("infrastruktur smtp gagal - konfigurasi environment variable tidak lengkap")
	}

	auth := smtp.PlainAuth("", smtpAuthEmail, smtpAuthPassword, smtpHost)

	subject := fmt.Sprintf("Pembayaran Berhasil - Pesanan #%s", orderID)
	body := fmt.Sprintf("Halo!\n\nKami ingin memberitahukan bahwa pembayaran untuk pesanan Anda dengan nomor %s telah berhasil diproses oleh sistem kami.\n\nTerima kasih telah berbelanja!\n", orderID)

	senderName := "MiniOrder Notif-Service"

	message := []byte(fmt.Sprintf("From: %s <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", senderName, smtpAuthEmail, email, subject, body))

	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, smtpAuthEmail, []string{email}, message)
	if err != nil {
		logger.ErrorContext(ctx, "koneksi SMTP Gmail gagal melempar pesan",
			slog.String("order_id", orderID),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to send real email: %w", err)
	}

	logger.InfoContext(ctx, "EMAIL BERHASIL TERBANG MELALUI GMAIL SMTP! 🚀",
		slog.String("order_id", orderID),
		slog.String("recipient_email", email),
	)

	return nil
}
