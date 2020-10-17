# SwissKnife

### SwissKnife as the name suggests is a collection of utilities that help to build services for the web written in GO

Split into small modules so that you dont have to import everything. There is a default implementation for most of the modules. Users can also create plugins by implementing the module interface.

```

go get -u github.com/adityak368/swissknife/<modulename>

```

### Crypto

- Basic Crypto functions for symmetric and asymmetric key encryptions

```go
    import "github.com/adityak368/swissknife/crypto"

    priv, pub := crypto.GenerateRsaKeyPair()

    if err := crypto.ExportRsaPrivateKeyToFile("privatekey.pem", priv); err != nil {
        return err
    }

    if err := crypto.ExportRsaPublicKeyToFile("pubkey.pub", pub); err != nil {
        return err
    }

    encryptedData, err := crypto.EncryptUsingSymmKey(msg, privKey)
    if err != nil {
        return err
    }

    data, err := crypto.DecryptUsingSymmKey(encryptedData, privKey)
    if err != nil {
        return err
    }

```

### Email

- Email Module for sending emails

```go
    import "github.com/adityak368/swissknife/email"
    import "github.com/adityak368/swissknife/knifemailer"

    mailer := knifemailer.New(email.MailerConfig{
        Host:     config.EmailHost,
        Port:     config.EmailPort,
        Username: config.EmailUsername,
        Password: config.EmailPassword,
    })
    mailer.StartDaemon()
    defer mailer.StopDaemon()
    mailer.SendMail(From, To, Subject, Body)

```

### Localization

- Localization module to extract locales and perform translations
- Supports translations in JSON format ( Ex: en.json, de.json ) with key value pairs

```go
    import "github.com/adityak368/swissknife/localization"
    import "github.com/adityak368/swissknife/localization/i18n"

    localizer := i18n.Localizer()
    localizer.LoadJSONLocalesFromFolder("res/locales") // locales contains file en-US.json
    translator := localization.Get().Translator("en-US")
    translated := translator.Tr("Key")
    //Or With params
    translated := translator.Tr("Key", "Params")

```

### Logger

- Logger Module for easy application logging
- Supports console and writing to file

```go
    import "github.com/adityak368/swissknife/logger"

    logger.InitConsoleLogger()
    defer logger.DestroyLogger()

    // Write to file
    logger.InitFileLogger("App.log")

    defer logger.DestroyLogger()

    logger.Debug.Println("Test")
    logger.Info.Println("Test")
    logger.Error.Println("Test")

```

### Middleware

- Middlewares for echo server
  - ErrorHandler
    - Handles Error Translations
  - Localization
    - Sets the Translator so that we can translate in the error handler, or any other part of our code
  - RateLimiter
    - RateLimits Requests
  - Tracing
    - Adds OpenTracing to our server

```go
    import "github.com/adityak368/swissknife/middleware/tracing"
    import "github.com/adityak368/swissknife/middleware/ratelimiter"
    import "github.com/adityak368/swissknife/middleware/localization"
    import "github.com/adityak368/swissknife/middleware/errorhandler"
    // Using Echo Middleware
    e := echo.New()

    e.Use(tracing.EchoTracingMiddleware(tracer, "AppName"))
    e.Use(localization.EchoLocalizer())
    e.Use(ratelimiter.RateLimitMiddleware())
    e.HTTPErrorHandler = errorhandler.EchoHTTPErrorHandlerMiddleware

```

### ObjectStore

- Supplies Helpers to Upload/Download File/Image to Amazon S3 Store

```go
    import "github.com/adityak368/swissknife/objectstore"
    import "github.com/adityak368/swissknife/objectstore/s3"

    store := s3.NewStore()

    // Using echo context
    avatar, err := c.FormFile("avatar")

    if err != nil {
        return err
    }

    src, err := avatar.Open()
    defer src.Close()
    if err != nil {
        return err
    }

    // Check for jpeg and png images
    bufReader := bufio.NewReader(src)
    buffer, err := bufReader.Peek(int(math.Min(512, float64(avatar.Size))))
    contentType := http.DetectContentType(buffer)
    if err != nil || !(contentType == "image/jpeg" || contentType == "image/png") {
        return err
    }

    user := c.Get("user").(*models.User)

    uploadResult, err := store.AddImage(config.S3AvatarBucket, user.ID.Hex(), bufReader)
    if err != nil {
        return err
    }

```

### Response

- Defines the input/output interfaces for the internal business handlers

```go
    import "github.com/adityak368/swissknife/response"

    func Login(c context.Context, user *User) (*response.Result, error) {
        if user.Authenticate() {
            return &response.Result{Data: map[string]string{token: "Token"}}, nil
        }
        return nil, response.NewError(http.StatusBadRequest, "InvalidEmailPassword")
    }

```

### Validation

- Provides helpers to perform input validation

```go
    import "github.com/adityak368/swissknife/validation/playground"

    type Test struct {
        A int `validate:"required,min=8,max=16"`
    }
    t := Test{A:1}
    validator := playground.New()
    validator.Validate(t)

```
