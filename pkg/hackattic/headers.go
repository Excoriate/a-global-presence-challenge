package hackattic

type CountryHeaderUserAgents struct {
	Country string
	Headers []string
}

func GetHeaders() []CountryHeaderUserAgents {
	// Set headers for 7 different countries.

	// Headers for US
	usHeaders := []string{
		"User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
		"X-Forwarded-For: 12.34.56.78", // An example US IP address
	}

	// Headers for UK
	ukHeaders := []string{
		"User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
		"X-Forwarded-For: 46.34.161.117", // An example UK IP address
	}

	// Headers for Germany
	germanyHeaders := []string{
		"User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
		"X-Forwarded-For: 77.182.74.143", // An example Germany IP address
	}

	// Headers for France
	franceHeaders := []string{
		"User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
		"X-Forwarded-For: 78.206.77.83", // An example France IP address
	}

	// Headers for Italy
	italyHeaders := []string{
		"User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
		"X-Forwarded-For: 95.233.211.222", // An example Italy IP address
	}

	// Headers for Spain
	spainHeaders := []string{
		"User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
		"X-Forwarded-For: 88.17.171.192", // An example Spain IP address
	}

	// Headers for Japan
	japanHeaders := []string{
		"User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
		"X-Forwarded-For: 58.84.43.69", // An example Japan IP address
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
