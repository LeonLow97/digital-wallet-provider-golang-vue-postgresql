package utils

const PasswordResetEmailBody = `
<html>
<head>
	<style>
		body {
			font-family: Arial, sans-serif;
			background-color: #f2f2f2;
			margin: 0;
			padding: 0;
		}
		.container {
			max-width: 600px;
			margin: 0 auto;
			padding: 20px;
			border: 1px solid #e0e0e0;
			border-radius: 5px;
			background-color: #ffffff;
		}
		h1 {
			color: #333;
			margin-bottom: 20px;
		}
		p {
			color: #666;
			margin-bottom: 15px;
		}
		.button-container {
			text-align: center;
		}
		.button {
			display: inline-block;
			background-color: #007bff;
			color: #fff;
			text-decoration: none;
			padding: 10px 20px;
			border-radius: 5px;
			margin-top: 10px;
		}
		.footer {
			margin-top: 30px;
			font-size: 12px;
			color: #999;
			text-align: center;
		}
	</style>
</head>
<body>
	<div class="container">
		<h1>Reset Your Password</h1>
		<p>Dear User,</p>
		<p>We have received a request to reset the password for your account. To proceed with resetting your password, please click the button below:</p>
		<div class="button-container">
			<a href="{{.ResetURL}}" class="button">Reset Password</a>
		</div>
		<p>If you did not request a password reset, you can safely ignore this email.</p>
		<div class="footer">
			<p>This email was sent automatically. Please do not reply to this email.</p>
			<p>Best regards,<br>Digital Wallet</p>
		</div>
	</div>
</body>
</html>
`
