Prom.ua SDK library for Go

## Install

```bash
go get github.com/freedayko/prom-golang-sdk
```

## Usage

```go
import "github.com/freedayko/prom-golang-sdk"
```

Get Products

```go
	p := prom.NewClient("<abcdef0123456789-you-prom-api-key>")

	var request = prom.ProductsRequest{
		GroupId: 0,
		Limit:   2000,
	}

	products, err := p.GetProducts(request)
	if err != nil {
		fmt.Printf("Error when requesting: %s", err)
		return
	}
```

Update Products

```go
	var editProducts []prom.ProductEdit
	
	editProduct := prom.NewProductEdit(product)
	editProduct.Discount.Value = 300

	editProducts = append(editProducts, editProduct)
	
	ids, err = p.UpdateProducts(editProducts)
```

