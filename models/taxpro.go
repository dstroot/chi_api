package models

import (
	"github.com/dstroot/chi_api/database"
)

// TaxPro is a Tax Professional
type TaxPro struct {
	EFIN           string `json:"efin"`
	CompanyName    string `json:"company_name"`
	ProductCount   int    `json:"product_count"`
	PremierPartner bool   `json:"premier_partner"`
}

// GetTaxpro returns a tax professional
func GetTaxpro(year string, efin string) ([]*TaxPro, error) {

	query := `
	SELECT TOP(1)
		E.EFIN,
		E.CompanyName,
		D.PriorVolume,
		PremierPartner = CASE WHEN D.PriorVolume < 250 THEN 0 ELSE 1 END
	FROM  eroyeardetail D
	RIGHT OUTER JOIN ero E on E.id = D.ero_id
	WHERE D.systemyear = ?
		AND D.status IN ('A', 'C', 'D')
		AND LastImportDate <> '' AND E.EFIN = ?;`

	rows, err := database.DB.Query(query, year, efin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*TaxPro, 0)
	pro := new(TaxPro)

	for rows.Next() {
		err1 := rows.Scan(&pro.EFIN, &pro.CompanyName, &pro.ProductCount, &pro.PremierPartner)
		if err1 != nil {
			return nil, err
		}
		results = append(results, pro)
	}
	if err2 := rows.Err(); err2 != nil {
		return nil, err2
	}
	return results, nil
}
