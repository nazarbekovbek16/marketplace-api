// repository/product_repository.go

package repository

import (
	"database/sql"
	"gorm.io/gorm"
	"marketplace-api/internal/models"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (pr *ProductRepository) CreateProduct(product *models.Product) error {
	if err := pr.db.Create(&product).Error; err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepository) UpdateProduct(product *models.Product) error {
	if product.Stock == 0 {
		if err := pr.db.Where("id = ?", product.ID).Model(product).Update("stock", 0).Error; err != nil {
			return err
		}
	}
	if err := pr.db.Where("id = ?", product.ID).Updates(&product).Error; err != nil {
		return err
	}

	if product.ImgURLs == nil {
		var empty []string
		return pr.db.Model(&product).Where("id = ?", product.ID).Update("img_urls", empty).Error
	}
	return nil
}

func (pr *ProductRepository) DeleteProduct(productID int64) error {
	if err := pr.db.Delete(&models.Product{}, productID).Error; err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepository) GetProductByID(productID int64) (*models.Product, error) {
	var product models.Product
	if err := pr.db.First(&product, productID).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (pr *ProductRepository) GetProductsByDistributorID(productName string, filters models.Filters, distributorID int64) ([]*models.Product, models.Metadata, error) {
	rows, err := pr.db.Table("products").Select("count(*) OVER()",
		"id", "category", "product_name", "product_description",
		"price", "img_urls", "minimum_quantity", "stock", "city").Where(
		"(to_tsvector('simple', product_name) @@ plainto_tsquery('simple', ?) OR ? = '') AND distributor_id = ?", productName, productName, distributorID).
		Order(filters.SortColumn() + " " + filters.SortDirection()).
		Order("id ASC").
		Limit(filters.Limit()).
		Limit(filters.Offset()).Rows()
	if err != nil {
		return nil, models.Metadata{}, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	totalRecords := 0
	var products []*models.Product

	for rows.Next() {
		var product models.Product

		err := rows.Scan(
			&totalRecords,
			&product.ID,
			&product.Category,
			&product.ProductName,
			&product.ProductDescription,
			&product.Price,
			&product.ImgURLs,
			&product.MinimumQuantity,
			&product.Stock,
			&product.City,
		)
		if err != nil {
			return nil, models.Metadata{}, err
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, models.Metadata{}, err
	}
	metadata := models.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return products, metadata, nil
}

func (pr *ProductRepository) GetProducts(productName string, filters models.Filters) ([]*models.Product, models.Metadata, error) {
	rows, err := pr.db.Table("products").Select("count(*) OVER()",
		"id", "category", "product_name", "product_description",
		"price", "img_urls", "minimum_quantity", "stock", "city").Where(
		"(products.stock != 0) AND (to_tsvector('simple', product_name) @@ plainto_tsquery('simple', ?) OR ? = '')", productName, productName).
		Order(filters.SortColumn() + " " + filters.SortDirection()).
		Order("id ASC").
		Limit(filters.Limit()).
		Limit(filters.Offset()).Rows()
	if err != nil {
		return nil, models.Metadata{}, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	totalRecords := 0
	var products []*models.Product

	for rows.Next() {
		var product models.Product

		err := rows.Scan(
			&totalRecords,
			&product.ID,
			&product.Category,
			&product.ProductName,
			&product.ProductDescription,
			&product.Price,
			&product.ImgURLs,
			&product.MinimumQuantity,
			&product.Stock,
			&product.City,
		)
		if err != nil {
			return nil, models.Metadata{}, err
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, models.Metadata{}, err
	}
	metadata := models.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return products, metadata, nil
}

func (pr *ProductRepository) GetEmail(userID int64) (string, error) {
	var user models.User
	err := pr.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

func (pr *ProductRepository) CreatReview(review *models.Review) error {
	if err := pr.db.Create(&review).Error; err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepository) GetReviews(id int64, role string) ([]models.Review, error) {
	var reviews []models.Review
	if role == "store" {
		if err := pr.db.Where("store_id = ?", id).Find(&reviews).Error; err != nil {
			return nil, err
		}
	} else if role == "distributor" {
		if err := pr.db.Where("distributor_id = ?", id).Find(&reviews).Error; err != nil {
			return nil, err
		}
	} else if role == "product" {
		if err := pr.db.Where("product_id = ?", id).Find(&reviews).Error; err != nil {
			return nil, err
		}
	}

	return reviews, nil
}

func (pr *ProductRepository) DeleteByReviewId(id int64) error {
	if err := pr.db.Delete(&models.Review{}, id).Error; err != nil {
		return err
	}
	return nil
}
