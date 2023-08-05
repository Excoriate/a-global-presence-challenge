package hackattic

type CountryHeaderUserAgents struct {
	Country string
	Headers []string
}

func GetHeaders() []CountryHeaderUserAgents {
	// Set headers for 7 different countries.

	// Headers for US
	usHeaders := []string{
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko)",
		"Chrome/56.0.2924.87 Safari/537.36",
	}

	// Headers for UK
	ukHeaders := []string{
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko)",
		"Chrome/56.0.2924.87 Safari/537.36",
	}

	// Headers for Germany
	germanyHeaders := []string{
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko)",
		"Chrome/56.0.2924.87 Safari/537.36",
	}

	// Headers for France
	franceHeaders := []string{
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko)",
		"Chrome/56.0.2924.87 Safari/537.36",
	}

	// Headers for Italy
	italyHeaders := []string{
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko)",
		"Chrome/56.0.2924.87 Safari/537.36",
	}

	// Headers for Spain
	spainHeaders := []string{
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko)",
		"Chrome/56.0.2924.87 Safari/537.36",
	}

	// Headers for Japan
	japanHeaders := []string{
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko)",
		"Chrome/56.0.2924.87 Safari/537.36",
	}

	return []CountryHeaderUserAgents{
		{
			Country: "US",
			Headers: usHeaders,
		},
		{
			Country: "UK",
			Headers: ukHeaders,
		},
		{
			Country: "Germany",
			Headers: germanyHeaders,
		},
		{
			Country: "France",
			Headers: franceHeaders,
		},
		{
			Country: "Italy",
			Headers: italyHeaders,
		},
		{
			Country: "Spain",
			Headers: spainHeaders,
		},
		{
			Country: "Japan",
			Headers: japanHeaders,
		},
	}
}
