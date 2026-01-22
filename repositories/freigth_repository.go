package repositories

import (
	"api_frete/database"
	"api_frete/interfaces"
	"api_frete/models"
	"database/sql"
	"fmt"
	"time"
)

type FreightRepository struct {
	db *sql.DB
}

func NewFreightRepository() interfaces.IFreightRepository {
	return &FreightRepository{
		db: database.Conn,
	}
}

func (r *FreightRepository) SaveQuote(carriers []models.CarrierInfo) error {
	if r.db == nil {
		return fmt.Errorf("conexão com banco não estabelecida")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
        INSERT INTO freight_carriers (carrier_name, service, deadline, price, created_at)
        VALUES ($1, $2, $3, $4, $5)		
    `)
	if err != nil {
		return fmt.Errorf("erro ao preparar statement: %v", err)
	}
	defer stmt.Close()

	for _, carrier := range carriers {
		_, err = stmt.Exec(
			carrier.Name,
			carrier.Service,
			carrier.Deadline,
			carrier.Price,
			time.Now(),
		)
		if err != nil {
			return fmt.Errorf("erro ao inserir transportadora: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao fazer commit: %v", err)
	}

	return nil
}

func (r *FreightRepository) GetCarrierStatistics(limit int) (*models.FreightStatisticsResponse, error) {
	if r.db == nil {
		return nil, fmt.Errorf("conexão com banco não estabelecida")
	}

	limitStr := fmt.Sprintf("limit %d", limit)
	sql := `select fc.carrier_name,
       			   sum(fc.price) as total,
                   avg(fc.price) as average,
                   count(1) as qty
             from (
                   select * 
                     from freight_carriers
                   order by created_at desc
				   %s
				) fc
             group by fc.carrier_name `

	if limit <= 0 {
		sql = fmt.Sprintf(sql, "")
	} else {
		sql = fmt.Sprintf(sql, limitStr)
	}

	rows, err := r.db.Query(sql)

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar estatísticas: %v", err)
	}
	defer rows.Close()

	carriers := []models.CarrierStatistics{}
	var minPrice, maxPrice float64

	for rows.Next() {
		var carrier models.CarrierStatistics
		err := rows.Scan(
			&carrier.Name,
			&carrier.TotalFreight,
			&carrier.AverageFreight,
			&carrier.QtyResults,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler dados: %v", err)
		}
		carriers = append(carriers, carrier)
	}

	err = r.db.QueryRow(`
        select coalesce(min(price), 0) as min_price,
               coalesce(max(price), 0) as max_price
          from freight_carriers
    `).Scan(&minPrice, &maxPrice)

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar min/max: %v", err)
	}

	return &models.FreightStatisticsResponse{
		Carrier:  carriers,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
	}, nil
}
