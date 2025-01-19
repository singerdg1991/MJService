package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	logger "github.com/hoitek/Logger"
	"github.com/hoitek/Logger/engines"
	loggerPort "github.com/hoitek/Logger/ports"
	"github.com/hoitek/Maja-Service/agenda"
	"github.com/hoitek/Maja-Service/agenda/tasks"
	"github.com/hoitek/Maja-Service/config"
	"github.com/hoitek/Maja-Service/database"
	"github.com/hoitek/Maja-Service/database/migrations"
	"github.com/hoitek/Maja-Service/eventstore"
	"github.com/hoitek/Maja-Service/internal/ability"
	"github.com/hoitek/Maja-Service/internal/address"
	"github.com/hoitek/Maja-Service/internal/ai"
	"github.com/hoitek/Maja-Service/internal/archive"
	"github.com/hoitek/Maja-Service/internal/city"
	"github.com/hoitek/Maja-Service/internal/contracttype"
	"github.com/hoitek/Maja-Service/internal/customer"
	"github.com/hoitek/Maja-Service/internal/cycle"
	"github.com/hoitek/Maja-Service/internal/diagnose"
	"github.com/hoitek/Maja-Service/internal/email"
	"github.com/hoitek/Maja-Service/internal/equipment"
	"github.com/hoitek/Maja-Service/internal/evaluation"
	"github.com/hoitek/Maja-Service/internal/geartype"
	"github.com/hoitek/Maja-Service/internal/healthcheck"
	"github.com/hoitek/Maja-Service/internal/keikkala"
	"github.com/hoitek/Maja-Service/internal/languageskill"
	"github.com/hoitek/Maja-Service/internal/license"
	"github.com/hoitek/Maja-Service/internal/limitation"
	"github.com/hoitek/Maja-Service/internal/medicine"
	"github.com/hoitek/Maja-Service/internal/notification"
	"github.com/hoitek/Maja-Service/internal/oauth2"
	"github.com/hoitek/Maja-Service/internal/otp"
	"github.com/hoitek/Maja-Service/internal/paymenttype"
	"github.com/hoitek/Maja-Service/internal/permission"
	"github.com/hoitek/Maja-Service/internal/prescription"
	"github.com/hoitek/Maja-Service/internal/punishment"
	"github.com/hoitek/Maja-Service/internal/push"
	"github.com/hoitek/Maja-Service/internal/quiz"
	"github.com/hoitek/Maja-Service/internal/report"
	"github.com/hoitek/Maja-Service/internal/reward"
	"github.com/hoitek/Maja-Service/internal/role"
	"github.com/hoitek/Maja-Service/internal/section"
	"github.com/hoitek/Maja-Service/internal/service"
	"github.com/hoitek/Maja-Service/internal/servicegrade"
	"github.com/hoitek/Maja-Service/internal/serviceoption"
	"github.com/hoitek/Maja-Service/internal/servicetype"
	"github.com/hoitek/Maja-Service/internal/shifttype"
	"github.com/hoitek/Maja-Service/internal/staff"
	"github.com/hoitek/Maja-Service/internal/staffclub/attention"
	"github.com/hoitek/Maja-Service/internal/staffclub/grace"
	"github.com/hoitek/Maja-Service/internal/staffclub/holiday"
	"github.com/hoitek/Maja-Service/internal/staffclub/warning"
	"github.com/hoitek/Maja-Service/internal/stafftype"
	"github.com/hoitek/Maja-Service/internal/static"
	"github.com/hoitek/Maja-Service/internal/ticket"
	"github.com/hoitek/Maja-Service/internal/todo"
	"github.com/hoitek/Maja-Service/internal/trash"
	"github.com/hoitek/Maja-Service/internal/user"
	"github.com/hoitek/Maja-Service/internal/vehicle"
	"github.com/hoitek/Maja-Service/internal/vehicletype"
	"github.com/hoitek/Maja-Service/internal/welcome"
	"github.com/hoitek/Maja-Service/messagebroker"
	"github.com/hoitek/Maja-Service/router"
	"github.com/hoitek/Maja-Service/storage"
	"github.com/hoitek/Maja-Service/tikka"
)

func main() {
	// Load configuration
	config.LoadDefault()

	// Initialize Logger
	logger.Default = logger.Initialize[loggerPort.LoggerEngineType](
		&engines.LoggerEngineStdout{},
	)

	// Initialize agenda
	agendaManager := agenda.Setup()

	// Migrate database
	migrations.MigrateDB()

	// Connect to Postgresql
	pDB := database.ConnectPostgresDB()

	// Connect to MongoDB
	database.ConnectMongoDB()

	// Connect to MinIO server
	var (
		MinIOEndpoint  = config.AppConfig.MinioEndpoint
		MinIOAccessKey = config.AppConfig.MinioAccessKey
		MinIOSecretKey = config.AppConfig.MinioSecretKey
	)
	minioStorage := storage.NewMinIO(MinIOEndpoint, MinIOAccessKey, MinIOSecretKey)
	_, err := minioStorage.Connect()
	if err != nil {
		log.Fatalf("MinIO client err: %v", err)
	}

	// Config multi language
	//bundle := transify.NewBundle(language.English)
	//bundle.SetTranslationDir(path.Join(config.GetRootPath(), "i18n", "translations"))
	//if err := bundle.LoadMessages(); err != nil {
	//	log.Fatal(err)
	//}

	// Connect to EventStore server (ES)
	es, err := eventstore.Setup(
		config.AppConfig.GrpcHost,
		config.AppConfig.GrpcPort,
		config.AppConfig.EventStoreRabbitMQHost,
		config.AppConfig.EventStoreRabbitMQPort,
		config.AppConfig.EventStoreRabbitMQUser,
		config.AppConfig.EventStoreRabbitMQPassword,
	)
	if err != nil {
		log.Fatal("Failed to setup event store", err)
	}

	// Connect to Tikka server
	_, err = tikka.Setup(
		config.AppConfig.OTPGRPCHost,
		config.AppConfig.OTPGRPCPort,
	)
	if err != nil {
		log.Fatal("Failed to setup tikka", err)
	}

	// Connect to RabbitMQ server
	mb, err := messagebroker.ConnectRabbitMQ(
		config.AppConfig.RabbitMQHost,
		config.AppConfig.RabbitMQPort,
		config.AppConfig.RabbitMQUser,
		config.AppConfig.RabbitMQPassword,
	)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ", err)
	}

	// Initialize router instance
	r := router.Init()

	// Set allow origin config
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"POST", "HEAD", "GET", "PATCH", "OPTIONS", "PUT", "DELETE"})
	headersOk := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	corsHandler := handlers.CORS(originsOk, methodsOk, headersOk)(r)

	// Register welcome handler
	welcome.Module.Setup(config.LoadWelcomeConfig()).RegisterHTTP(r)

	// Register healthcheck handler
	healthcheck.Module.Setup(config.LoadHealthCheckConfig()).RegisterHTTP(r)

	// Register ai handler
	ai.Module.Setup(config.LoadAIConfig()).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register city handler
	city.Module.Setup(config.LoadCityConfig()).
		SetDatabase(database.PostgresDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register gearType handler
	geartype.Module.Setup(config.LoadGearTypeConfig()).
		SetDatabase(database.PostgresDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register equipment handler
	equipment.Module.Setup(config.LoadEquipmentConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register reward handler
	reward.Module.Setup(config.LoadRewardConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register punishment handler
	punishment.Module.Setup(config.LoadPunishmentConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register languageSkill handler
	languageskill.Module.Setup(config.LoadLanguageSkillConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register staffType handler
	stafftype.Module.Setup(config.LoadStaffTypeConfig()).
		SetDatabase(database.PostgresDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register permission handler
	permission.Module.Setup(config.LoadPermissionConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register role handler
	role.Module.Setup(config.LoadRoleConfig()).
		SetDatabase(database.PostgresDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register shiftType handler
	shifttype.Module.Setup(config.LoadShiftTypeConfig()).
		SetDatabase(database.PostgresDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register contractType handler
	contracttype.Module.Setup(config.LoadContractTypeConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register section handler
	section.Module.Setup(config.LoadSectionConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register trash handler
	trash.Module.Setup(config.LoadTrashConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register vehicleType handler
	vehicletype.Module.Setup(config.LoadVehicleTypeConfig()).
		SetDatabase(database.PostgresDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register ability handler
	ability.Module.Setup(config.LoadAbilityConfig()).
		SetDatabase(database.PostgresDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register paymentType handler
	paymenttype.Module.Setup(config.LoadPaymentTypeConfig()).
		SetDatabase(database.PostgresDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register user handler
	user.Module.Setup(config.LoadUserConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterWorkers().
		RegisterHTTP(r)

	// Register vehicle handler
	vehicle.Module.Setup(config.LoadVehicleConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register staff handler
	staff.Module.Setup(config.LoadStaffConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		//ForceMigrateAndSeed().
		RegisterWorkers().
		RegisterHTTP(r)

	// Register medicine handler
	medicine.Module.Setup(config.LoadMedicineConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register prescription handler
	prescription.Module.Setup(config.LoadPrescriptionConfig()).
		SetDatabase(database.PostgresDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register diagnose handler
	diagnose.Module.Setup(config.LoadDiagnoseConfig()).
		SetDatabase(database.PostgresDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register serviceType handler
	servicetype.Module.Setup(config.LoadServiceTypeConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register limitation handler
	limitation.Module.Setup(config.LoadLimitationConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register license handler
	license.Module.Setup(config.LoadLicenseConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register serviceGrade handler
	servicegrade.Module.Setup(config.LoadServiceGradeConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register service handler
	service.Module.Setup(config.LoadServiceConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register address handler
	address.Module.Setup(config.LoadAddressConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register customer handler
	customer.Module.Setup(config.LoadCustomerConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterWorkers().
		RegisterHTTP(r)

	// Register cycle handler
	cycleModule := cycle.Module.Setup(config.LoadCycleConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		SetMessageBroker(mb).
		RegisterQueues().
		RegisterCron()
	cycleModule.RegisterHTTP(r)

	// Register archive handler
	archive.Module.Setup(config.LoadArchiveConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register report category
	servicetype.Module.Setup(config.LoadServiceTypeConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register report
	serviceoption.Module.Setup(config.LoadServiceOptionConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register staffClub grace
	grace.Module.Setup(config.LoadGraceConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register staffClub attention
	attention.Module.Setup(config.LoadAttentionConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register staffClub warning
	warning.Module.Setup(config.LoadWarningConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register staffClub holiday
	holiday.Module.Setup(config.LoadHolidayConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register otp handler
	otp.Module.Setup(config.LoadOTPConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register oauth2 handler
	oauth2.Module.Setup(config.LoadOAuth2Config()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register ticket handler
	ticket.Module.Setup(config.LoadTicketConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register todos handler
	todo.Module.Setup(config.LoadTodoConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register notification handler
	notification.Module.Setup(config.LoadNotificationConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register quiz handler
	quiz.Module.Setup(config.LoadQuizConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register email handler
	email.Module.Setup(config.LoadEmailConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register keikkala handler
	keikkala.Module.Setup(config.LoadKeikkalaConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register evaluation handler
	evaluation.Module.Setup(config.LoadEvaluationConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register push handler
	push.Module.Setup(config.LoadPushConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register report handler
	report.Module.Setup(config.LoadReportConfig()).
		SetDatabases(database.PostgresDB, database.MongoDB).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register static handler
	static.Module.Setup(config.LoadStaticConfig()).
		SetMinIOStorage(minioStorage).
		RegisterHTTP(r)

	// Register reminder agenda task
	reminderTask := agenda.NewTask[time.Duration]("Every 1 Day", tasks.AgendaTaskReminder, 24*time.Hour)
	agendaManager.AddTask(reminderTask)

	// Initialize migrator
	_, err = migrations.Apply(pDB)
	if err != nil {
		log.Fatalf("failed to initialize migrator: %v", err)
	}

	// Retrieve host name and port from config
	var (
		PROTOCOL        = config.AppConfig.Protocol
		HOST_ADDRESS    = config.AppConfig.HostAddress
		HOST_URI        = config.AppConfig.HostUri
		PORT            = config.AppConfig.Port
		SWAGGER_UI_PATH = "/apidocs"
	)

	// Create address and uri
	ADDR := fmt.Sprintf("%s:%d", HOST_ADDRESS, PORT)
	ADDR_URI := fmt.Sprintf("%s:%d", HOST_URI, PORT)

	// Run agenda manager
	agendaManager.Run()

	// Start http server
	msg := make(chan error)
	go func() {
		log.Printf("Server started at %s://%s\n", PROTOCOL, ADDR_URI)
		log.Printf("Swagger UI accessible at %s://%s\n", PROTOCOL, ADDR_URI+SWAGGER_UI_PATH)
		if err := http.ListenAndServe(ADDR, corsHandler); err != nil {
			msg <- err
		}
	}()

	// Listen for interrupt signal
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		msg <- fmt.Errorf("%s", <-c)
	}()

	// Print error message
	log.Println(<-msg)

	// Close gRPC connection
	es.Close()

	// Close message broker
	mb.Close()

	// Exit Application
	os.Exit(1)
}
