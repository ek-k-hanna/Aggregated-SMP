module hannaekthesis/main

go 1.16

replace hannaekthesis/vahss_bprp => ./vahss_bprp

replace hannaekthesis/bulletproof => ./bulletproofs_go

replace hannaekthesis/vahss_ccs => ./vahss_ccs

replace hannaekthesis/ccs08 => ./ccs08_go

replace hannaekthesis/vahss => ./vahss

require (
	hannaekthesis/bulletproof v0.0.0-00010101000000-000000000000 // indirect
	hannaekthesis/ccs08 v0.0.0-00010101000000-000000000000 // indirect
	hannaekthesis/vahss v0.0.0-00010101000000-000000000000 // indirect
	hannaekthesis/vahss_bprp v0.0.0-00010101000000-000000000000 // indirect
	hannaekthesis/vahss_ccs v0.0.0-00010101000000-000000000000 // indirect
)
