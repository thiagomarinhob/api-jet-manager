package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Modelos (repetidos aqui para evitar dependências cíclicas)
// Definições dos tipos
// Tipos de enumeração
type UserType string
type TableStatus string
type OrderStatus string
type OrderType string
type ProductType string
type TransactionType string
type TransactionCategory string
type SubscriptionStatus string

// Constantes para os tipos
const (
	// UserType
	UserTypeSuperAdmin UserType = "superadmin"
	UserTypeAdmin      UserType = "admin"
	UserTypeManager    UserType = "manager"
	UserTypeStaff      UserType = "staff"

	// TableStatus
	TableStatusFree     TableStatus = "free"
	TableStatusOccupied TableStatus = "occupied"
	TableStatusReserved TableStatus = "reserved"

	// OrderStatus
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusReady     OrderStatus = "ready"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"

	// OrderType
	OrderTypeInHouse  OrderType = "in_house"
	OrderTypeDelivery OrderType = "delivery"
	OrderTypeTakeaway OrderType = "takeaway"

	// ProductCategory
	ProductCategoryFood    ProductType = "food"
	ProductCategoryDrink   ProductType = "drink"
	ProductCategoryDessert ProductType = "dessert"

	// TransactionType
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"

	// TransactionCategory
	TransactionCategorySales       TransactionCategory = "sales"
	TransactionCategoryOther       TransactionCategory = "other_income"
	TransactionCategoryIngredients TransactionCategory = "ingredients"
	TransactionCategoryUtilities   TransactionCategory = "utilities"
	TransactionCategorySalaries    TransactionCategory = "salaries"
	TransactionCategoryRent        TransactionCategory = "rent"
	TransactionCategoryEquipment   TransactionCategory = "equipment"
	TransactionCategoryMaintenance TransactionCategory = "maintenance"

	// SubscriptionStatus
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusInactive SubscriptionStatus = "inactive"
	SubscriptionStatusTrial    SubscriptionStatus = "trial"
)

// Modelos
type Restaurant struct {
	ID               uuid.UUID          `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name             string             `gorm:"size:100;not null"`
	Description      string             `gorm:"size:255"`
	Address          string             `gorm:"size:255"`
	Phone            string             `gorm:"size:20"`
	Email            string             `gorm:"size:100"`
	Logo             string             `gorm:"size:255"`
	SubscriptionPlan string             `gorm:"size:50"`
	Status           SubscriptionStatus `gorm:"size:20;not null;default:'trial'"`
	TrialEndsAt      *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type User struct {
	ID           uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name         string     `gorm:"size:100;not null"`
	Email        string     `gorm:"size:100;uniqueIndex;not null"`
	Password     string     `gorm:"size:100;not null"`
	Type         UserType   `gorm:"size:20;not null;default:'staff'"`
	RestaurantID *uuid.UUID `gorm:"type:uuid"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Table struct {
	ID             uuid.UUID   `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	RestaurantID   uuid.UUID   `gorm:"type:uuid;not null"`
	Number         int         `gorm:"not null"`
	Capacity       int         `gorm:"not null"`
	Status         TableStatus `gorm:"size:20;not null;default:'free'"`
	CurrentOrderID *uuid.UUID  `gorm:"type:uuid"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ProductCategory struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	RestaurantID uuid.UUID `gorm:"type:uuid;not null"`
	Name         string    `gorm:"size:100;not null"`
	Description  string    `gorm:"size:255"`
	Active       bool      `gorm:"default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Product struct {
	ID              uuid.UUID        `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	RestaurantID    uuid.UUID        `gorm:"type:uuid;not null"`
	Name            string           `gorm:"size:100;not null"`
	Description     string           `gorm:"size:255"`
	Price           float64          `gorm:"not null"`
	CategoryID      uuid.UUID        `gorm:"type:uuid;not null"`
	ProductCategory *ProductCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Type            ProductType      `gorm:"size:20;not null"`
	InStock         bool             `gorm:"default:true"`
	ImageURL        string           `gorm:"size:255"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Order struct {
	ID              uuid.UUID   `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	RestaurantID    uuid.UUID   `gorm:"type:uuid;not null"`
	TableID         *uuid.UUID  `gorm:"type:uuid"`
	UserID          uuid.UUID   `gorm:"type:uuid;not null"`
	CustomerName    string      `gorm:"size:100"`
	CustomerPhone   string      `gorm:"size:20"`
	CustomerEmail   string      `gorm:"size:100"`
	Type            OrderType   `gorm:"size:20;not null;default:'in_house'"`
	Status          OrderStatus `gorm:"size:20;not null;default:'pending'"`
	TotalAmount     float64     `gorm:"not null;default:0"`
	Notes           string      `gorm:"size:255"`
	DeliveryAddress string      `gorm:"size:255"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	PaidAt          *time.Time
	DeliveredAt     *time.Time
}

type OrderItem struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	Quantity  int       `gorm:"not null;default:1"`
	Price     float64   `gorm:"not null"`
	Notes     string    `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FinancialTransaction struct {
	ID            uuid.UUID           `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	RestaurantID  uuid.UUID           `gorm:"type:uuid;not null"`
	Type          TransactionType     `gorm:"size:20;not null"`
	Category      TransactionCategory `gorm:"size:30;not null"`
	Amount        float64             `gorm:"not null"`
	Description   string              `gorm:"size:255"`
	OrderID       *uuid.UUID          `gorm:"type:uuid"`
	UserID        uuid.UUID           `gorm:"type:uuid;not null"`
	PaymentMethod string              `gorm:"size:30"`
	Date          time.Time           `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Dados de exemplo para seed (mesmos de antes)
var (
	// Nomes fictícios de restaurantes
	restaurantNames = []string{
		"Sabor Brasileiro", "Cantina Italiana", "Sushi Express",
		"Churrascaria Gaúcha", "Veggie Delights", "Pizzaria Napolitana",
		"Burger House", "Taco Fiesta", "Mediterranean Delight", "Asian Fusion",
	}

	// ... (resto das variáveis igual ao seu script original)
	// Descrições de restaurantes
	restaurantDescriptions = []string{
		"Autêntica culinária brasileira com atmosfera acolhedora",
		"Sabores tradicionais da Itália no coração da cidade",
		"Sushi fresco e pratos japoneses preparados por chefs experientes",
		"Carnes nobres grelhadas no estilo gaúcho tradicional",
		"Restaurante vegetariano com ingredientes orgânicos e frescos",
		"Pizzas artesanais assadas em forno a lenha",
		"Hambúrgueres gourmet com ingredientes premium",
		"Autêntica comida mexicana com receitas tradicionais",
		"Sabores do Mediterrâneo com azeites importados",
		"Fusão de sabores asiáticos com toque contemporâneo",
	}

	// Categorias padrão para produtos
	defaultCategories = []struct {
		Name        string
		Description string
	}{
		{"Pratos Principais", "Pratos principais do cardápio"},
		{"Entradas", "Entradas e aperitivos"},
		{"Bebidas Não Alcoólicas", "Refrigerantes, sucos e água"},
		{"Bebidas Alcoólicas", "Cervejas, vinhos e drinks"},
		{"Sobremesas", "Doces e sobremesas"},
	}

	// Categorias de alimentos
	foodItems = []struct {
		Name        string
		Description string
		Type        ProductType
		Price       float64
	}{
		{"Feijoada Completa", "Tradicional prato brasileiro com arroz, farofa e couve", ProductCategoryFood, 49.90},
		{"Filé Mignon", "Corte nobre grelhado com molho de vinho", ProductCategoryFood, 75.90},
		{"Sushi Combo", "18 peças variadas de sushi e sashimi", ProductCategoryFood, 89.90},
		{"Pizza Margherita", "Molho de tomate, muçarela e manjericão", ProductCategoryFood, 45.90},
		{"Salada Caesar", "Alface romana, croutons, parmesão e molho especial", ProductCategoryFood, 32.90},
		{"Hambúrguer Gourmet", "Blend especial com queijo cheddar e bacon", ProductCategoryFood, 39.90},
		{"Pescada Grelhada", "Peixe fresco grelhado com ervas", ProductCategoryFood, 59.90},
		{"Risoto de Funghi", "Arroz arbóreo com cogumelos frescos e parmesão", ProductCategoryFood, 55.90},
		{"Lasanha à Bolonhesa", "Camadas de massa com molho de carne e bechamel", ProductCategoryFood, 47.90},
		{"Frango Parmegiana", "Peito de frango empanado com molho de tomate e queijo gratinado", ProductCategoryFood, 49.90},

		{"Água Mineral", "Garrafa 500ml", ProductCategoryDrink, 6.90},
		{"Refrigerante", "Lata 350ml", ProductCategoryDrink, 7.90},
		{"Suco Natural", "300ml de suco fresco de laranja, abacaxi ou limão", ProductCategoryDrink, 12.90},
		{"Caipirinha", "Tradicional drink brasileiro de limão", ProductCategoryDrink, 22.90},
		{"Vinho Tinto", "Taça da casa 150ml", ProductCategoryDrink, 25.90},
		{"Cerveja Artesanal", "Garrafa 355ml", ProductCategoryDrink, 18.90},
		{"Café Espresso", "50ml de café forte", ProductCategoryDrink, 7.90},
		{"Chá Gelado", "300ml com limão e hortelã", ProductCategoryDrink, 11.90},

		{"Pudim de Leite", "Clássico pudim com calda de caramelo", ProductCategoryDessert, 16.90},
		{"Mousse de Chocolate", "Cremoso com raspas de chocolate belga", ProductCategoryDessert, 19.90},
		{"Tiramisu", "Sobremesa italiana com café e mascarpone", ProductCategoryDessert, 24.90},
		{"Cheesecake", "Com calda de frutas vermelhas", ProductCategoryDessert, 22.90},
		{"Sorvete", "Duas bolas com calda a escolha", ProductCategoryDessert, 18.90},
		{"Petit Gateau", "Bolo de chocolate com centro derretido e sorvete", ProductCategoryDessert, 25.90},
	}

	// Nomes e emails de funcionários
	staffNames = []string{
		"João Silva", "Maria Oliveira", "Carlos Santos", "Ana Pereira",
		"Pedro Costa", "Lúcia Fernandes", "Roberto Almeida", "Juliana Lima",
		"Fernando Ribeiro", "Patrícia Souza", "Marcos Gomes", "Camila Rodrigues",
		"Ricardo Martins", "Luiza Carvalho", "Eduardo Ferreira", "Beatriz Nunes",
	}

	staffEmails = []string{
		"joao.silva@exemplo.com", "maria.oliveira@exemplo.com", "carlos.santos@exemplo.com",
		"ana.pereira@exemplo.com", "pedro.costa@exemplo.com", "lucia.fernandes@exemplo.com",
		"roberto.almeida@exemplo.com", "juliana.lima@exemplo.com", "fernando.ribeiro@exemplo.com",
		"patricia.souza@exemplo.com", "marcos.gomes@exemplo.com", "camila.rodrigues@exemplo.com",
		"ricardo.martins@exemplo.com", "luiza.carvalho@exemplo.com", "eduardo.ferreira@exemplo.com",
		"beatriz.nunes@exemplo.com",
	}

	// Planos de assinatura
	subscriptionPlans = []string{"basic", "standard", "premium"}

	// Status de assinatura
	subscriptionStatuses = []SubscriptionStatus{
		SubscriptionStatusActive, SubscriptionStatusTrial, SubscriptionStatusInactive,
	}

	// Tipos de usuário
	userTypes = []UserType{UserTypeAdmin, UserTypeManager, UserTypeStaff}

	// Formas de pagamento
	paymentMethods = []string{
		"Dinheiro", "Cartão de Crédito", "Cartão de Débito", "PIX", "Vale Refeição",
	}

	// Endereços
	addresses = []string{
		"Av. Paulista, 1234 - São Paulo, SP",
		"Rua Augusta, 500 - São Paulo, SP",
		"Av. Atlântica, 900 - Rio de Janeiro, RJ",
		"Av. Boa Viagem, 1100 - Recife, PE",
		"Rua das Flores, 321 - Curitiba, PR",
		"Av. Beira Mar, 1500 - Fortaleza, CE",
		"Rua Padre Chagas, 250 - Porto Alegre, RS",
		"Av. Getúlio Vargas, 800 - Belo Horizonte, MG",
		"Rua das Palmeiras, 100 - Salvador, BA",
		"Av. Dom Pedro I, 2000 - Manaus, AM",
	}

	// Nomes de clientes
	customerNames = []string{
		"Rodrigo Mendes", "Sandra Vieira", "Thiago Alves", "Paula Dias",
		"Felipe Cardoso", "Renata Moraes", "Bruno Ferreira", "Daniela Castro",
		"Gabriel Souza", "Cristina Pires", "Leandro Machado", "Isabela Costa",
		"Hugo Almeida", "Tatiana Martins", "Leonardo Santos", "Amanda Rodrigues",
	}

	// Emails de clientes
	customerEmails = []string{
		"rodrigo.mendes@email.com", "sandra.vieira@email.com", "thiago.alves@email.com",
		"paula.dias@email.com", "felipe.cardoso@email.com", "renata.moraes@email.com",
		"bruno.ferreira@email.com", "daniela.castro@email.com", "gabriel.souza@email.com",
		"cristina.pires@email.com", "leandro.machado@email.com", "isabela.costa@email.com",
		"hugo.almeida@email.com", "tatiana.martins@email.com", "leonardo.santos@email.com",
		"amanda.rodrigues@email.com",
	}

	// Telefones de clientes
	customerPhones = []string{
		"(11) 98765-4321", "(11) 99876-5432", "(21) 98765-1234", "(21) 99876-2345",
		"(81) 98765-3456", "(81) 99876-4567", "(41) 98765-5678", "(41) 99876-6789",
		"(85) 98765-7890", "(85) 99876-8901", "(51) 98765-9012", "(51) 99876-0123",
		"(31) 98765-1234", "(31) 99876-2345", "(71) 98765-3456", "(71) 99876-4567",
	}

	// Notas para pedidos
	orderNotes = []string{
		"Sem cebola, por favor", "Ponto da carne: mal passado", "Molho à parte",
		"Sem pimenta", "Bebidas bem geladas", "Sem glúten", "Sem lactose",
		"Extra queijo", "Sem tomate", "Ponto da carne: bem passado",
		"", "", "", "", // Algumas entradas vazias para variar
	}

	// Categorias de despesas
	expenseCategories = []TransactionCategory{
		TransactionCategoryIngredients, TransactionCategoryUtilities,
		TransactionCategorySalaries, TransactionCategoryRent,
		TransactionCategoryEquipment, TransactionCategoryMaintenance,
	}

	// Descrições de transações financeiras
	transactionDescriptions = map[TransactionCategory][]string{
		TransactionCategorySales: {
			"Venda do dia", "Receita de delivery", "Reserva privada", "Evento corporativo",
		},
		TransactionCategoryOther: {
			"Aluguel de espaço", "Venda de equipamento usado", "Patrocínio", "Royalties",
		},
		TransactionCategoryIngredients: {
			"Compra de carnes", "Compra de vegetais", "Compra de bebidas", "Compra de grãos",
			"Compra de laticínios", "Compra de frutos do mar", "Compra de especiarias",
		},
		TransactionCategoryUtilities: {
			"Conta de luz", "Conta de água", "Conta de gás", "Internet", "Telefone",
		},
		TransactionCategorySalaries: {
			"Folha de pagamento", "Bônus", "Horas extras", "Benefícios",
		},
		TransactionCategoryRent: {
			"Aluguel mensal", "Taxa de condomínio", "IPTU",
		},
		TransactionCategoryEquipment: {
			"Compra de forno", "Compra de geladeira", "Compra de utensílios",
			"Compra de móveis", "Sistema de POS",
		},
		TransactionCategoryMaintenance: {
			"Manutenção de equipamentos", "Serviço de limpeza", "Reparo hidráulico",
			"Reparo elétrico", "Pintura", "Dedetização",
		},
	}
)

func generateUUID() uuid.UUID {
	return uuid.New()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Carrega variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente")
	}

	// Obtem os dados do banco de dados
	dbHost := "localhost"
	dbPort := getEnv("BLUEPRINT_DB_PORT", "5432")
	dbUser := getEnv("BLUEPRINT_DB_USERNAME", "docker")
	dbPassword := getEnv("BLUEPRINT_DB_PASSWORD", "docker")
	dbName := getEnv("BLUEPRINT_DB_DATABASE", "jetmanager")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	// Constrói a string de conexão
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	// Conecta ao banco de dados
	log.Println("Conectando ao banco de dados...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}
	log.Println("Conexão estabelecida com sucesso!")

	// Verifica se as tabelas existem
	if !db.Migrator().HasTable(&Restaurant{}) {
		log.Println("Tabelas não existem. Execute as migrações primeiro.")
		return
	}

	// Inicia a transação para o seed
	tx := db.Begin()

	// Cria superadmin se não existir
	var superadminCount int64
	tx.Model(&User{}).Where("type = ?", UserTypeSuperAdmin).Count(&superadminCount)
	if superadminCount == 0 {
		createSuperAdmin(tx)
	}

	// Seed de restaurantes e relacionados
	createRestaurants(tx)

	// Commit da transação
	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Erro ao commitar a transação: %v", err)
	}

	log.Println("Seed concluído com sucesso!")
}

func createSuperAdmin(db *gorm.DB) {
	log.Println("Criando superadmin...")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("superadmin"), bcrypt.DefaultCost)
	superadmin := User{
		ID:       generateUUID(),
		Name:     "Super Admin",
		Email:    "superadmin@example.com",
		Password: string(hashedPassword),
		Type:     UserTypeSuperAdmin,
	}
	db.Create(&superadmin)
}

func createRestaurants(db *gorm.DB) {
	log.Println("Criando restaurantes e dados relacionados...")

	// Criar entre 5 e 10 restaurantes rand.Intn(6) + 5
	numRestaurants := 1
	for i := 0; i < numRestaurants; i++ {
		// Cria restaurante
		restaurant := createRestaurant(db, i)

		// Cria funcionários para o restaurante
		users := createUsers(db, restaurant.ID)

		fmt.Println("USERS == ", users)

		categories := createProductCategories(db, restaurant.ID)

		fmt.Println("Categories == ", categories)
		// Cria produtos (menu)
		products := createProducts(db, restaurant.ID, categories)
		fmt.Println("Products == ", products)
		// Cria mesas
		tables := createTables(db, restaurant.ID)
		fmt.Println("Tables == ", tables)
		// Cria pedidos
		orders := createOrders(db, restaurant.ID, users, tables, products)
		fmt.Println(orders)

		// Cria transações financeiras
		// createFinancialTransactions(db, restaurant.ID, users, orders)
	}
}

func createRestaurant(db *gorm.DB, index int) Restaurant {
	statusIndex := rand.Intn(len(subscriptionStatuses))
	planIndex := rand.Intn(len(subscriptionPlans))

	now := time.Now()
	var trialEndsAt *time.Time

	if subscriptionStatuses[statusIndex] == SubscriptionStatusTrial {
		endDate := now.AddDate(0, 1, 0) // 1 mês de teste
		trialEndsAt = &endDate
	}

	restaurant := Restaurant{
		ID:          generateUUID(),
		Name:        restaurantNames[index%len(restaurantNames)],
		Description: restaurantDescriptions[index%len(restaurantDescriptions)],
		Address:     addresses[index%len(addresses)],
		Phone: fmt.Sprintf("(1%d) 9%d%d%d%d-%d%d%d%d",
			rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9),
			rand.Intn(9), rand.Intn(9), rand.Intn(9), rand.Intn(9)),
		Email:            fmt.Sprintf("contato@%s.com", sanitizeForEmail(restaurantNames[index%len(restaurantNames)])),
		Logo:             fmt.Sprintf("https://example.com/logos/%d.png", index+1),
		SubscriptionPlan: subscriptionPlans[planIndex],
		Status:           subscriptionStatuses[statusIndex],
		TrialEndsAt:      trialEndsAt,
		CreatedAt:        time.Now().AddDate(0, -rand.Intn(6), -rand.Intn(30)), // Entre 0 e 6 meses atrás
	}
	restaurant.UpdatedAt = restaurant.CreatedAt

	db.Create(&restaurant)
	return restaurant
}

func createUsers(db *gorm.DB, restaurantID uuid.UUID) []User {
	var users []User

	// Cria um admin
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := User{
		ID:           generateUUID(),
		Name:         staffNames[0],
		Email:        fmt.Sprintf("admin@%s.example.com", restaurantID),
		Password:     string(adminPassword),
		Type:         UserTypeAdmin,
		RestaurantID: &restaurantID,
		CreatedAt:    time.Now().AddDate(0, -rand.Intn(3), 0), // Entre 0 e 3 meses atrás
	}
	admin.UpdatedAt = admin.CreatedAt
	db.Create(&admin)
	users = append(users, admin)

	// Cria entre 1 e 3 gerentes
	numManagers := rand.Intn(3) + 1
	for i := 0; i < numManagers; i++ {
		managerPassword, _ := bcrypt.GenerateFromPassword([]byte("manager123"), bcrypt.DefaultCost)
		staffIndex := (i + 1) % len(staffNames)
		manager := User{
			ID:           generateUUID(),
			Name:         staffNames[staffIndex],
			Email:        fmt.Sprintf("manager%d@%s.example.com", i+1, restaurantID),
			Password:     string(managerPassword),
			Type:         UserTypeManager,
			RestaurantID: &restaurantID,
			CreatedAt:    time.Now().AddDate(0, -rand.Intn(3), 0), // Entre 0 e 3 meses atrás
		}
		manager.UpdatedAt = manager.CreatedAt
		db.Create(&manager)
		users = append(users, manager)
	}

	// Cria entre 3 e 10 funcionários
	numStaff := rand.Intn(8) + 3
	for i := 0; i < numStaff; i++ {
		staffPassword, _ := bcrypt.GenerateFromPassword([]byte("staff123"), bcrypt.DefaultCost)
		staffIndex := (i + 1 + numManagers) % len(staffNames)
		staff := User{
			ID:           generateUUID(),
			Name:         staffNames[staffIndex],
			Email:        fmt.Sprintf("staff%d@%s.example.com", i+1, restaurantID),
			Password:     string(staffPassword),
			Type:         UserTypeStaff,
			RestaurantID: &restaurantID,
			CreatedAt:    time.Now().AddDate(0, -rand.Intn(2), 0), // Entre 0 e 2 meses atrás
		}
		staff.UpdatedAt = staff.CreatedAt
		db.Create(&staff)
		users = append(users, staff)
	}

	return users
}

// NOVA FUNÇÃO: Cria categorias de produtos para um restaurante
func createProductCategories(db *gorm.DB, restaurantID uuid.UUID) []ProductCategory {
	log.Printf("Criando categorias de produtos para o restaurante %s...", restaurantID)

	var productCategory []ProductCategory

	// Cria as categorias padrão
	for _, catInfo := range defaultCategories {
		category := ProductCategory{
			ID:           generateUUID(),
			RestaurantID: restaurantID,
			Name:         catInfo.Name,
			Description:  catInfo.Description,
			Active:       true,
			CreatedAt:    time.Now().AddDate(0, -rand.Intn(3), 0), // Entre 0 e 3 meses atrás
		}
		category.UpdatedAt = category.CreatedAt

		db.Create(&category)
		productCategory = append(productCategory, category)
	}

	log.Printf("Criadas %d categorias para o restaurante %s", len(productCategory), restaurantID)
	return productCategory
}

func createProducts(db *gorm.DB, restaurantID uuid.UUID, productCategories []ProductCategory) []Product {
	var products []Product

	// Cria seleção de produtos do menu
	numProducts := rand.Intn(10) + 15 // Entre 15 e 25 produtos
	for i := 0; i < numProducts; i++ {
		foodIndex := i % len(foodItems)
		food := foodItems[foodIndex]

		// Varia o preço um pouco
		priceVariation := rand.Float64()*0.2 - 0.1 // Entre -10% e +10%
		price := food.Price * (1 + priceVariation)

		var indexC = rand.Intn(len(productCategories))

		product := Product{
			ID:              generateUUID(),
			RestaurantID:    restaurantID,
			Name:            food.Name,
			Description:     food.Description,
			Price:           float64(int(price*100)) / 100, // Arredonda para 2 casas decimais
			CategoryID:      productCategories[indexC].ID,
			ProductCategory: &productCategories[indexC],
			Type:            food.Type,
			InStock:         rand.Float64() < 0.9, // 90% dos produtos estão em estoque
			ImageURL:        fmt.Sprintf("https://example.com/products/%d.jpg", i+1),
			CreatedAt:       time.Now().AddDate(0, -rand.Intn(3), 0), // Entre 0 e 3 meses atrás
		}
		product.UpdatedAt = product.CreatedAt
		db.Create(&product)
		products = append(products, product)
	}

	return products
}

func createTables(db *gorm.DB, restaurantID uuid.UUID) []Table {
	var tables []Table

	// Cria entre 5 e 20 mesas
	numTables := rand.Intn(16) + 5
	for i := 0; i < numTables; i++ {
		// Define capacidade (entre 2 e 8)
		capacity := (rand.Intn(4) + 1) * 2

		// Status da mesa (maioria livre)
		statusRoll := rand.Float64()
		var status TableStatus
		if statusRoll < 0.7 {
			status = TableStatusFree
		} else if statusRoll < 0.9 {
			status = TableStatusOccupied
		} else {
			status = TableStatusReserved
		}

		table := Table{
			ID:           generateUUID(),
			RestaurantID: restaurantID,
			Number:       i + 1,
			Capacity:     capacity,
			Status:       status,
			CreatedAt:    time.Now().AddDate(0, -rand.Intn(3), 0), // Entre 0 e 3 meses atrás
		}
		table.UpdatedAt = table.CreatedAt
		db.Create(&table)
		tables = append(tables, table)
	}

	return tables
}

func createOrders(db *gorm.DB, restaurantID uuid.UUID, users []User, tables []Table, products []Product) []Order {
	var orders []Order

	// Cria entre 20 e 50 pedidos
	// numOrders := rand.Intn(31) + 20
	var numOrders = 5
	// Variedade de datas para pedidos (nos últimos 90 dias)
	now := time.Now()

	for i := 0; i < numOrders; i++ {
		// Seleciona usuário aleatório (que criou o pedido)
		user := users[rand.Intn(len(users))]

		// Define tipo de pedido
		orderTypeRoll := rand.Float64()
		var orderType OrderType
		var tableID *uuid.UUID

		if orderTypeRoll < 0.6 { // 60% são pedidos no local
			orderType = OrderTypeInHouse
			// Associa a uma mesa aleatória
			table := tables[rand.Intn(len(tables))]
			tableID = &table.ID
		} else if orderTypeRoll < 0.9 { // 30% são delivery
			orderType = OrderTypeDelivery
		} else { // 10% são takeaway
			orderType = OrderTypeTakeaway
		}

		// Define status do pedido
		statusRoll := rand.Float64()
		var orderStatus OrderStatus
		var paidAt *time.Time
		var deliveredAt *time.Time

		// Para pedidos mais antigos, maior chance de estarem concluídos
		daysPast := rand.Intn(90) // Entre 0 e 90 dias atrás
		orderDate := now.AddDate(0, 0, -daysPast)

		// Conforme o pedido é mais antigo, maior a chance de estar completo
		completionFactor := float64(daysPast) / 90.0 // 0 a 1

		if statusRoll < 0.1*(1-completionFactor) { // Pedidos mais recentes têm maior chance de estarem pendentes
			orderStatus = OrderStatusPending
		} else if statusRoll < 0.2*(1-completionFactor) {
			orderStatus = OrderStatusPreparing
		} else if statusRoll < 0.3*(1-completionFactor) {
			orderStatus = OrderStatusReady
		} else if statusRoll < 0.4*(1-completionFactor) {
			orderStatus = OrderStatusDelivered
			deliveryTime := orderDate.Add(time.Duration(30+rand.Intn(30)) * time.Minute)
			deliveredAt = &deliveryTime
		} else if statusRoll < 0.95 { // Maioria dos pedidos mais antigos estão pagos
			orderStatus = OrderStatusPaid
			paymentTime := orderDate.Add(time.Duration(45+rand.Intn(45)) * time.Minute)
			paidAt = &paymentTime

			// Se for delivery, também foi entregue
			if orderType == OrderTypeDelivery {
				deliveryTime := orderDate.Add(time.Duration(30+rand.Intn(30)) * time.Minute)
				deliveredAt = &deliveryTime
			}
		} else { // 5% foram cancelados
			orderStatus = OrderStatusCancelled
		}

		// Dados do cliente
		var customerName, customerEmail, customerPhone, deliveryAddress string

		// Para delivery e takeaway, preenche dados do cliente
		if orderType != OrderTypeInHouse {
			customerIndex := rand.Intn(len(customerNames))
			customerName = customerNames[customerIndex]
			customerEmail = customerEmails[customerIndex]
			customerPhone = customerPhones[customerIndex]

			// Para delivery, adiciona endereço
			if orderType == OrderTypeDelivery {
				deliveryAddress = addresses[rand.Intn(len(addresses))]
			}
		}

		// Cria o pedido
		order := Order{
			ID:              generateUUID(),
			RestaurantID:    restaurantID,
			TableID:         tableID,
			UserID:          user.ID,
			CustomerName:    customerName,
			CustomerPhone:   customerPhone,
			CustomerEmail:   customerEmail,
			Type:            orderType,
			Status:          orderStatus,
			Notes:           orderNotes[rand.Intn(len(orderNotes))],
			DeliveryAddress: deliveryAddress,
			CreatedAt:       time.Now(),
			UpdatedAt:       orderDate,
			PaidAt:          paidAt,
			DeliveredAt:     deliveredAt,
		}

		db.Create(&order)

		// Adiciona itens ao pedido
		totalAmount := createOrderItems(db, order.ID, products)

		// Atualiza o valor total do pedido
		db.Model(&order).Update("total_amount", totalAmount)

		// Atualiza o order no slice local
		order.TotalAmount = totalAmount
		orders = append(orders, order)

		// Se o pedido está ocupando uma mesa (em andamento e no local), atualiza status da mesa
		if orderType == OrderTypeInHouse && (orderStatus == OrderStatusPending ||
			orderStatus == OrderStatusPreparing || orderStatus == OrderStatusReady) {

			db.Model(&Table{}).Where("id = ?", tableID).Updates(map[string]interface{}{
				"status":        TableStatusOccupied,
				"current_order": order.ID,
			})
		}
	}

	return orders
}

func createOrderItems(db *gorm.DB, orderID uuid.UUID, products []Product) float64 {
	var totalAmount float64

	// Adiciona entre 1 e 6 itens ao pedido
	numItems := rand.Intn(6) + 1

	// Evita produtos duplicados rastreando os já selecionados
	selectedProducts := make(map[string]bool)

	for i := 0; i < numItems; i++ {
		// Seleciona um produto aleatório (sem repetir)
		var product Product
		for {
			product = products[rand.Intn(len(products))]
			if !selectedProducts[product.ID.String()] {
				selectedProducts[product.ID.String()] = true
				break
			}

			// Se já temos muitos produtos selecionados, quebra o loop
			if len(selectedProducts) >= len(products) {
				break
			}
		}

		// Define quantidade (1 a 3)
		quantity := rand.Intn(3) + 1

		// Cria o item do pedido
		item := OrderItem{
			ID:        generateUUID(),
			OrderID:   orderID,
			ProductID: product.ID,
			Quantity:  quantity,
			Price:     product.Price,
			Notes:     orderNotes[rand.Intn(len(orderNotes))],
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		db.Create(&item)

		// Soma ao valor total
		totalAmount += product.Price * float64(quantity)
	}

	return totalAmount
}

func createFinancialTransactions(db *gorm.DB, restaurantID uuid.UUID, users []User, orders []Order) {
	log.Printf("Criando transações financeiras para o restaurante %s...", restaurantID)

	// Administrador para transações (ou gerente)
	var adminUser User
	for _, user := range users {
		if user.Type == UserTypeAdmin || user.Type == UserTypeManager {
			adminUser = user
			break
		}
	}

	// 1. Cria transações para cada pedido pago
	for _, order := range orders {
		if order.Status == OrderStatusPaid && order.PaidAt != nil {
			// Método de pagamento aleatório
			paymentMethod := paymentMethods[rand.Intn(len(paymentMethods))]

			// Cria transação de receita
			transaction := FinancialTransaction{
				ID:            generateUUID(),
				RestaurantID:  restaurantID,
				Type:          TransactionTypeIncome,
				Category:      TransactionCategorySales,
				Amount:        order.TotalAmount,
				Description:   fmt.Sprintf("Venda #%d - %s", order.ID, paymentMethod),
				OrderID:       &order.ID,
				UserID:        adminUser.ID,
				PaymentMethod: paymentMethod,
				Date:          *order.PaidAt,
				CreatedAt:     *order.PaidAt,
				UpdatedAt:     *order.PaidAt,
			}

			db.Create(&transaction)
		}
	}

	// 2. Cria despesas mensais para os últimos 3 meses
	now := time.Now()

	// Para cada mês (últimos 3)
	for month := 0; month < 3; month++ {
		// Data do mês atual
		currentMonth := now.AddDate(0, -month, 0)

		// Cria entre 4 e 10 despesas por mês
		numExpenses := rand.Intn(7) + 4

		for i := 0; i < numExpenses; i++ {
			// Seleciona categoria de despesa aleatória
			category := expenseCategories[rand.Intn(len(expenseCategories))]
			descriptions := transactionDescriptions[category]

			// Valor da despesa (entre R$50 e R$5000, dependendo da categoria)
			var amount float64
			switch category {
			case TransactionCategoryIngredients:
				amount = 500 + rand.Float64()*1500
			case TransactionCategoryUtilities:
				amount = 100 + rand.Float64()*400
			case TransactionCategorySalaries:
				amount = 1000 + rand.Float64()*4000
			case TransactionCategoryRent:
				amount = 2000 + rand.Float64()*3000
			case TransactionCategoryEquipment:
				amount = 300 + rand.Float64()*2000
			case TransactionCategoryMaintenance:
				amount = 100 + rand.Float64()*500
			}

			// Arredonda para 2 casas decimais
			amount = float64(int(amount*100)) / 100

			// Data da transação (dia aleatório do mês)
			day := rand.Intn(28) + 1 // Entre 1 e 28
			transactionDate := time.Date(
				currentMonth.Year(), currentMonth.Month(), day,
				10+rand.Intn(8), rand.Intn(60), rand.Intn(60), 0,
				time.Local,
			)

			// Cria a transação
			transaction := FinancialTransaction{
				ID:            generateUUID(),
				RestaurantID:  restaurantID,
				Type:          TransactionTypeExpense,
				Category:      category,
				Amount:        amount,
				Description:   descriptions[rand.Intn(len(descriptions))],
				UserID:        adminUser.ID,
				PaymentMethod: paymentMethods[rand.Intn(len(paymentMethods))],
				Date:          transactionDate,
				CreatedAt:     transactionDate,
				UpdatedAt:     transactionDate,
			}

			db.Create(&transaction)
		}
	}
}

// sanitizeForEmail remove caracteres especiais para criar um email válido
func sanitizeForEmail(name string) string {
	// Converte para minúsculo
	result := strings.ToLower(name)

	// Remove espaços e caracteres especiais
	result = strings.ReplaceAll(result, " ", "")
	result = strings.ReplaceAll(result, "á", "a")
	result = strings.ReplaceAll(result, "à", "a")
	result = strings.ReplaceAll(result, "ã", "a")
	result = strings.ReplaceAll(result, "â", "a")
	result = strings.ReplaceAll(result, "é", "e")
	result = strings.ReplaceAll(result, "ê", "e")
	result = strings.ReplaceAll(result, "í", "i")
	result = strings.ReplaceAll(result, "ó", "o")
	result = strings.ReplaceAll(result, "ô", "o")
	result = strings.ReplaceAll(result, "õ", "o")
	result = strings.ReplaceAll(result, "ú", "u")
	result = strings.ReplaceAll(result, "ç", "c")

	return result
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
