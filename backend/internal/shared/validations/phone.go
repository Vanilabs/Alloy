package validations

import (
	"fmt"
	"strings"

	"github.com/biter777/countries"
	"github.com/go-playground/validator/v10"
)

var (
	callingCodeMap = initCallingCodeMap()
	phoneLengthMap = initPhoneLengthMap()
)

// phoneLengthByCountryCode maps country calling codes to expected phone number lengths (without country code)
// These are the typical lengths for local phone numbers in each country
var phoneLengthByCountryCode = map[string]int{
	"+234": 10, // Nigeria
	"+1":   10, // USA/Canada
	"+44":  10, // UK
	"+33":  9,  // France
	"+49":  11, // Germany
	"+91":  10, // India
	"+86":  11, // China
	"+81":  10, // Japan
	"+7":   10, // Russia
	"+55":  11, // Brazil
	"+27":  9,  // South Africa
	"+61":  9,  // Australia
	"+34":  9,  // Spain
	"+39":  10, // Italy
	"+46":  9,  // Sweden
	"+47":  8,  // Norway
	"+31":  9,  // Netherlands
	"+32":  9,  // Belgium
	"+41":  9,  // Switzerland
	"+43":  10, // Austria
	"+45":  8,  // Denmark
	"+358": 9,  // Finland
	"+353": 9,  // Ireland
	"+351": 9,  // Portugal
	"+30":  10, // Greece
	"+48":  9,  // Poland
	"+420": 9,  // Czech Republic
	"+36":  9,  // Hungary
	"+40":  10, // Romania
	"+359": 9,  // Bulgaria
	"+385": 9,  // Croatia
	"+381": 9,  // Serbia
	"+382": 8,  // Montenegro
	"+383": 8,  // Kosovo
	"+386": 8,  // Slovenia
	"+387": 8,  // Bosnia
	"+389": 8,  // North Macedonia
	"+355": 9,  // Albania
	"+356": 8,  // Malta
	"+357": 8,  // Cyprus
	"+370": 8,  // Lithuania
	"+371": 8,  // Latvia
	"+372": 8,  // Estonia
	"+352": 9,  // Luxembourg
	"+377": 8,  // Monaco
	"+378": 8,  // San Marino
	"+212": 9,  // Morocco
	"+213": 9,  // Algeria
	"+216": 8,  // Tunisia
	"+20":  10, // Egypt
	"+233": 9,  // Ghana
	"+225": 8,  // Côte d'Ivoire
	"+221": 9,  // Senegal
	"+220": 7,  // Gambia
	"+224": 8,  // Guinea
	"+226": 8,  // Burkina Faso
	"+227": 8,  // Niger
	"+228": 8,  // Togo
	"+229": 8,  // Benin
	"+230": 7,  // Mauritius
	"+231": 7,  // Liberia
	"+232": 8,  // Sierra Leone
	"+235": 8,  // Chad
	"+236": 8,  // Central African Republic
	"+237": 9,  // Cameroon
	"+238": 7,  // Cape Verde
	"+239": 7,  // São Tomé and Príncipe
	"+240": 9,  // Equatorial Guinea
	"+241": 8,  // Gabon
	"+242": 9,  // Republic of the Congo
	"+243": 9,  // DR Congo
	"+244": 9,  // Angola
	"+245": 7,  // Guinea-Bissau
	"+246": 7,  // British Indian Ocean Territory
	"+247": 4,  // Ascension Island
	"+248": 7,  // Seychelles
	"+249": 9,  // Sudan
	"+250": 9,  // Rwanda
	"+251": 9,  // Ethiopia
	"+252": 8,  // Somalia
	"+253": 8,  // Djibouti
	"+257": 8,  // Burundi
	"+258": 9,  // Mozambique
	"+260": 9,  // Zambia
	"+261": 9,  // Madagascar
	"+262": 9,  // Réunion/Mayotte
	"+263": 9,  // Zimbabwe
	"+264": 9,  // Namibia
	"+265": 9,  // Malawi
	"+266": 8,  // Lesotho
	"+267": 8,  // Botswana
	"+268": 8,  // Eswatini
	"+269": 7,  // Comoros
	"+290": 4,  // Saint Helena
	"+291": 7,  // Eritrea
	"+297": 7,  // Aruba
	"+298": 6,  // Faroe Islands
	"+299": 6,  // Greenland
	"+350": 8,  // Gibraltar
	"+354": 7,  // Iceland
	"+373": 8,  // Moldova
	"+374": 8,  // Armenia
	"+375": 9,  // Belarus
	"+380": 9,  // Ukraine
	"+421": 9,  // Slovakia
	"+422": 0,  // (unassigned)
	"+423": 7,  // Liechtenstein
	"+500": 5,  // Falkland Islands
	"+501": 7,  // Belize
	"+502": 8,  // Guatemala
	"+503": 8,  // El Salvador
	"+504": 8,  // Honduras
	"+505": 8,  // Nicaragua
	"+506": 8,  // Costa Rica
	"+507": 8,  // Panama
	"+508": 6,  // Saint Pierre and Miquelon
	"+509": 8,  // Haiti
	"+590": 9,  // Guadeloupe
	"+591": 8,  // Bolivia
	"+592": 7,  // Guyana
	"+593": 9,  // Ecuador
	"+594": 9,  // French Guiana
	"+595": 9,  // Paraguay
	"+596": 9,  // Martinique
	"+597": 7,  // Suriname
	"+598": 8,  // Uruguay
	"+599": 7,  // Curaçao
	"+670": 8,  // East Timor
	"+672": 5,  // Australian External Territories
	"+673": 7,  // Brunei
	"+674": 7,  // Nauru
	"+675": 8,  // Papua New Guinea
	"+676": 7,  // Tonga
	"+677": 7,  // Solomon Islands
	"+678": 7,  // Vanuatu
	"+679": 7,  // Fiji
	"+680": 7,  // Palau
	"+681": 6,  // Wallis and Futuna
	"+682": 5,  // Cook Islands
	"+683": 4,  // Niue
	"+684": 0,  // (unassigned)
	"+685": 5,  // Samoa
	"+686": 5,  // Kiribati
	"+687": 6,  // New Caledonia
	"+688": 5,  // Tuvalu
	"+689": 6,  // French Polynesia
	"+690": 4,  // Tokelau
	"+691": 7,  // Micronesia
	"+692": 7,  // Marshall Islands
	"+850": 8,  // North Korea
	"+852": 8,  // Hong Kong
	"+853": 8,  // Macau
	"+855": 9,  // Cambodia
	"+856": 9,  // Laos
	"+880": 10, // Bangladesh
	"+886": 9,  // Taiwan
	"+960": 7,  // Maldives
	"+961": 8,  // Lebanon
	"+962": 9,  // Jordan
	"+963": 9,  // Syria
	"+964": 10, // Iraq
	"+965": 8,  // Kuwait
	"+966": 9,  // Saudi Arabia
	"+967": 9,  // Yemen
	"+968": 8,  // Oman
	"+970": 9,  // Palestine
	"+971": 9,  // UAE
	"+972": 9,  // Israel
	"+973": 8,  // Bahrain
	"+974": 8,  // Qatar
	"+975": 8,  // Bhutan
	"+976": 8,  // Mongolia
	"+977": 10, // Nepal
	"+992": 9,  // Tajikistan
	"+993": 8,  // Turkmenistan
	"+994": 9,  // Azerbaijan
	"+995": 9,  // Georgia
	"+996": 9,  // Kyrgyzstan
	"+998": 9,  // Uzbekistan
}

// initCallingCodeMap creates a map of ITU-T E.164 calling codes (with and without +) for fast lookup
func initCallingCodeMap() map[string]bool {
	m := make(map[string]bool)
	allCountries := countries.All()
	for _, c := range allCountries {
		callCodes := c.CallCodes()
		for _, callCode := range callCodes {
			// Store both with and without + prefix for flexibility
			codeStr := callCode.String()
			if codeStr != "" {
				// Store with + prefix
				m[codeStr] = true
				// Store without + prefix
				if strings.HasPrefix(codeStr, "+") {
					m[codeStr[1:]] = true
				} else {
					m["+"+codeStr] = true
				}
			}
		}
	}
	return m
}

// initPhoneLengthMap creates a map of country codes to phone number lengths
// This includes both the predefined mapping and dynamically builds from countries package
func initPhoneLengthMap() map[string]int {
	m := make(map[string]int)

	// Copy predefined mappings
	for code, length := range phoneLengthByCountryCode {
		m[code] = length
		// Also store without + prefix
		if strings.HasPrefix(code, "+") {
			m[code[1:]] = length
		} else {
			m["+"+code] = length
		}
	}

	// For countries not in the predefined map, use a default range
	// Most countries have phone numbers between 7-11 digits
	allCountries := countries.All()
	for _, c := range allCountries {
		callCodes := c.CallCodes()
		for _, callCode := range callCodes {
			codeStr := callCode.String()
			if codeStr != "" {
				// If not already in map, use default of 10 digits
				if _, exists := m[codeStr]; !exists {
					m[codeStr] = 10 // Default length
					if strings.HasPrefix(codeStr, "+") {
						m[codeStr[1:]] = 10
					} else {
						m["+"+codeStr] = 10
					}
				}
			}
		}
	}

	return m
}

// GetPhoneNumberLength returns the expected phone number length for a given country code
func GetPhoneNumberLength(countryCode string) (int, bool) {
	// Normalize country code
	normalized := strings.TrimSpace(countryCode)
	if !strings.HasPrefix(normalized, "+") && len(normalized) > 0 {
		allDigits := true
		for _, r := range normalized {
			if r < '0' || r > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			normalized = "+" + normalized
		}
	}

	length, exists := phoneLengthMap[normalized]
	if !exists && strings.HasPrefix(normalized, "+") {
		length, exists = phoneLengthMap[normalized[1:]]
	}

	return length, exists
}

// validateCountryCode validates that the country code is a valid ITU-T E.164 calling code (e.g., +234, +1, 234, 1)
func validateCountryCode(fl validator.FieldLevel) bool {
	countryCode := strings.TrimSpace(fl.Field().String())

	// Empty country code is allowed (handled by omitempty or required_if)
	if countryCode == "" {
		return true
	}

	// Normalize: ensure it starts with + or is just digits
	normalized := countryCode
	if !strings.HasPrefix(normalized, "+") && len(normalized) > 0 {
		// If it doesn't start with +, check if it's all digits
		allDigits := true
		for _, r := range normalized {
			if r < '0' || r > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			normalized = "+" + normalized
		}
	}

	// Check if it's a valid ITU-T E.164 calling code using cached map
	_, exists := callingCodeMap[normalized]
	if !exists {
		// Also check without + prefix
		if strings.HasPrefix(normalized, "+") {
			_, exists = callingCodeMap[normalized[1:]]
		}
	}
	return exists
}

// ValidateCountryCode validates a phone country code string (ITU-T E.164 calling code) using the countries package
func ValidateCountryCode(countryCode string) error {
	countryCode = strings.TrimSpace(countryCode)

	if countryCode == "" {
		return &ValidationError{
			Field:   "country_code",
			Message: "Country code is required",
		}
	}

	// Normalize: ensure it starts with + or is just digits
	normalized := countryCode
	if !strings.HasPrefix(normalized, "+") && len(normalized) > 0 {
		// If it doesn't start with +, check if it's all digits
		allDigits := true
		for _, r := range normalized {
			if r < '0' || r > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			normalized = "+" + normalized
		}
	}

	// Check if it's a valid ITU-T E.164 calling code using cached map
	_, exists := callingCodeMap[normalized]
	if !exists {
		// Also check without + prefix
		if strings.HasPrefix(normalized, "+") {
			_, exists = callingCodeMap[normalized[1:]]
		}
	}

	if !exists {
		return &ValidationError{
			Field:   "country_code",
			Message: "Invalid country code. Please provide a valid ITU-T E.164 calling code (e.g., +234, +1, +44, +33)",
		}
	}

	return nil
}

// ValidatePhoneNumber validates that a phone number contains only digits and has the correct length for the country
func ValidatePhoneNumber(phoneNumber, countryCode string) error {
	phoneNumber = strings.TrimSpace(phoneNumber)

	if phoneNumber == "" {
		return &ValidationError{
			Field:   "phone_number",
			Message: "Phone number is required",
		}
	}

	// Check that phone number contains only digits
	for _, char := range phoneNumber {
		if char < '0' || char > '9' {
			return &ValidationError{
				Field:   "phone_number",
				Message: "Phone number must contain only digits",
			}
		}
	}

	// Validate country code first
	if countryCode == "" {
		return &ValidationError{
			Field:   "country_code",
			Message: "Country code is required when using phone number",
		}
	}

	if err := ValidateCountryCode(countryCode); err != nil {
		return err
	}

	// Get expected length for this country code
	expectedLength, exists := GetPhoneNumberLength(countryCode)
	if !exists {
		// If country code is valid but we don't have length info, use a default range
		if len(phoneNumber) < 7 || len(phoneNumber) > 15 {
			return &ValidationError{
				Field:   "phone_number",
				Message: "Phone number length is invalid for the provided country code",
			}
		}
		return nil
	}

	// Validate exact length
	if len(phoneNumber) != expectedLength {
		return &ValidationError{
			Field:   "phone_number",
			Message: fmt.Sprintf("Phone number must be exactly %d digits for country code %s", expectedLength, countryCode),
		}
	}

	return nil
}

// FormatPhoneToE164 formats a phone number to E.164 format by appending country code
// Returns the E.164 formatted phone number (e.g., +2348123456789)
func FormatPhoneToE164(phoneNumber, countryCode string) (string, error) {
	phoneNumber = strings.TrimSpace(phoneNumber)
	countryCode = strings.TrimSpace(countryCode)

	if phoneNumber == "" {
		return "", fmt.Errorf("phone number is required")
	}

	if countryCode == "" {
		return "", fmt.Errorf("country code is required")
	}

	// Normalize country code - ensure it has + prefix
	normalizedCode := countryCode
	if !strings.HasPrefix(normalizedCode, "+") {
		normalizedCode = "+" + normalizedCode
	}

	// Simply append country code to phone number to create E.164 format
	// Phone number is already validated to contain only digits and correct length
	e164 := normalizedCode + phoneNumber
	return e164, nil
}
