package handler

import (
	jsoncode "encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"soulstreet/internal/model"
	"soulstreet/internal/service"
	"soulstreet/pkg/json"
	"strconv"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)
	name := r.FormValue("name")
	priceStr := r.FormValue("price")
	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		json.SendJsonError(w, http.StatusBadRequest, errors.New("Erro ao parsear 'price'"))
		return
	}
	files := r.MultipartForm.File["images"]
	var imagesPaths []string
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			json.SendJsonError(w, http.StatusInternalServerError, errors.New("Erro ao processar as imagens"))
			return
		}
		defer file.Close()
		newFileName := uuid.New().String() + filepath.Ext(fileHeader.Filename)
		savePath := "uploads/" + newFileName

		outFile, err := os.Create(savePath)
		if err != nil {
			json.SendJsonError(w, http.StatusInternalServerError, errors.New("Erro ao salvar imagem"))
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			json.SendJsonError(w, http.StatusInternalServerError, errors.New("Erro ao escrever imagem"))
			return
		}
		imgDir := "/images/" + newFileName
		imagesPaths = append(imagesPaths, imgDir)
	}
	imagesPathsJson, err := jsoncode.Marshal(imagesPaths)
	if err != nil {
		json.SendJsonError(w, http.StatusInternalServerError, errors.New("Erro ao converter array de imagem em json"))
		return
	}

	sizesStr := r.FormValue("sizes")
	sizes, err := jsoncode.Marshal(sizesStr)
	if err != nil{
		json.SendJsonError(w, http.StatusInternalServerError, errors.New("Erro ao converter array de tamanhos em json"))
		return
	}

	product := model.Product{
		Name:   name,
		Price:  float32(price),
		Images: string(imagesPathsJson),
		Sizes:  string(sizes),
	}

	err = h.productService.CreateProduct(&product)
	if err != nil {
		json.SendJsonError(w, http.StatusInternalServerError, errors.New("Erro ao salvar produto"))
		fmt.Println(err)
		return
	}
	fmt.Println(product)
	json.SendJson(w, http.StatusCreated, product)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		json.SendJsonError(w, http.StatusBadRequest, errors.New("Query 'id' em branco"))
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		json.SendJsonError(w, http.StatusBadRequest, errors.New("Erro ao parsear 'id'"))
		return
	}
	product, err := h.productService.GetProductByID(int(id))
	if err != nil {
		json.SendJsonError(w, http.StatusBadRequest, fmt.Errorf("error: %v", err))
		return
	}
	json.SendJson(w, http.StatusOK, product)

}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.GetAll()
	if err != nil {
		json.SendJsonError(w, http.StatusInternalServerError, err)
		return
	}
	json.SendJson(w, http.StatusOK, products)
}

func (h *ProductHandler) GetProductByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		json.SendJsonError(w, http.StatusBadRequest, errors.New("Query 'name' em branco"))
		return
	}
	product, err := h.productService.GetProductByName(name)
	if err != nil {
		json.SendJsonError(w, http.StatusBadRequest, fmt.Errorf("error: %v", err))
		return
	}
	json.SendJson(w, http.StatusOK, product)
}
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {	
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		json.SendJsonError(w, http.StatusBadRequest, errors.New("Query 'id' em branco"))
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		json.SendJsonError(w, http.StatusBadRequest, errors.New("Erro ao parsear 'id'"))
		return
	}
	err = h.productService.DeleteProduct(int(id))
	if err != nil {
		json.SendJsonError(w, http.StatusBadRequest, fmt.Errorf("error: %v", err))
		return
	}
	json.SendJson(w, http.StatusOK, "Produto deletado com sucesso")
}