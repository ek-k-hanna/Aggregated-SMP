module hannaekthesis/vahss_bprp

go 1.16

replace hannaekthesis/bulletproof => ../bulletproofs_go

replace hannaekthesis/vahss => ../vahss

replace hannaekthesis/p256 => ../crypto/p256

replace hannaekthesis/util => ../util

require (
	github.com/ing-bank/zkrp v0.0.0-20200519071134-97a3cddb5627 // indirect
	hannaekthesis/bulletproof v0.0.0-00010101000000-000000000000 // indirect
	hannaekthesis/p256 v0.0.0-00010101000000-000000000000 // indirect
	hannaekthesis/util v0.0.0-00010101000000-000000000000 // indirect
	hannaekthesis/vahss v0.0.0-00010101000000-000000000000 // indirect
)
