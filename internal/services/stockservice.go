// Package services
package services

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/anderson89marques/roxb3/internal/adapters/drivers/rest/handlers/schemas"
	domain "github.com/anderson89marques/roxb3/internal/core/domain"
	"github.com/anderson89marques/roxb3/internal/core/ports/database"
	"github.com/anderson89marques/roxb3/internal/infra/config"
)

type StockService struct {
	repo      database.Repository
	conf      *config.Config
	stockChan chan []*domain.Stock
}

func NewStockService(repo database.Repository, conf *config.Config) *StockService {
	return &StockService{
		repo: repo,
		conf: conf,
	}
}

const RootPath = "/app/files"

func B3Files() ([]string, error) {
	// Returns all files to be processed.
	entries, err := os.ReadDir(RootPath)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, v := range entries {
		if v.IsDir() {
			continue
		}
		files = append(files, filepath.Join(RootPath, v.Name()))
	}
	return files, nil
}

func (s *StockService) IngestData() error {
	// Read the b3 directory
	// Process the files in optimed way
	// call repository for bulk insert
	start := time.Now()
	files, err := B3Files()
	if err != nil {
		return err
	}
	fmt.Printf("Files to be processed: %+v \n", files)

	// Create worker pool
	s.stockChan = make(chan []*domain.Stock, s.conf.NumWorkers)
	var wg sync.WaitGroup
	for i := 0; i < s.conf.NumWorkers; i++ {
		wg.Go(func() {
			fmt.Printf("Worker: %d started\n", i)
			for batch := range s.stockChan {
				if err := s.InsertBatch(batch); err != nil {
					fmt.Printf("Error inserting batch: %v\n", err)
				}
			}
		})
	}

	for _, fileName := range files {
		err = s.ProcessStockFile(fileName)
		if err != nil {
			panic(err)
		}
	}

	go func() {
		fmt.Println("Refreshing materialized view.")
		if err := s.repo.RefreshStockSummaryMeterializedView(); err != nil {
			fmt.Println("Error to refresh materialized view.")
		}
	}()

	close(s.stockChan)
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Execution took %s\n", elapsed)
	return nil
}

func (s *StockService) ProcessStockFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'

	_, err = reader.Read() // skip header
	if err != nil {
		return err
	}

	cont := 0
	var batch []*domain.Stock
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break // End of file
		}
		if err != nil {
			return err
		}

		stock, err := domain.ToStockModel(row)
		if err != nil {
			return err
		}
		cont += 1
		batch = append(batch, stock)

		if len(batch) >= s.conf.BatchSize {
			s.stockChan <- batch
			batch = nil
		}
	}
	// Send remaining batch
	if len(batch) > 0 {
		s.stockChan <- batch
	}

	fmt.Printf("Total: %d\n", cont)
	return nil
}

func (s *StockService) InsertBatch(batch []*domain.Stock) error {
	// fmt.Printf("InsertBatch Size: %d\n", len(batch))
	if err := s.repo.BulkInsert(batch); err != nil {
		return err
	}
	return nil
}

func (s *StockService) Search(queryParams *schemas.StockSchema) (*domain.StockSummary, error) {
	stockSummary, err := s.repo.Search(queryParams)
	if err != nil {
		return nil, err
	}
	return stockSummary, nil
}
