# SwissKnife

### SwissKnife as the name suggests is a collection of utilities that help to build services for the web written in GO

Split into small modules so that you dont have to import everything. There is a default implementation for most of the modules. Users can also create plugins by implementing the module interface.

```
go get -u github.com/adityak368/swissknife/<modulename>@main
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

    encryptedData, err := crypto.EncryptUsingSymmKey(data, privKey)
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
    import (
        "github.com/adityak368/swissknife/email"
        "github.com/adityak368/swissknife/email/knifemailer"
    )

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
    import (
        "github.com/adityak368/swissknife/localization"
        "github.com/adityak368/swissknife/localization/i18n"
    )

    localizer := i18n.Localizer()
    localizer.LoadJSONLocalesFromFolder("res/locales") // locales contains file en-US.json
    translator := localizer.Translator("en-US")
    translated := translator.Tr("Key")
    //Or With params
    translated := translator.Tr("Key", "Params")
```

### Logger

- Logger Module for easy application logging
- Supports console and writing to file

```go
    import "github.com/adityak368/swissknife/logger"

    // Write to file
    logger.SetLogLevel(logger.LogLevelInfo)
    logger.SetShowCallerInfo(true)
    logger.SetLogOutputFile("App.log")
    // logger.SetOutput(writer io.Writer)

    logger.Debug("Debug")
    logger.Debugf("Hello %s", "World")
    logger.Info("Info")
    logger.Warn("Warn")
    logger.Error("Error")
    logger.Critical("Critical")
    logger.Trace("Trace")
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
    import (
        "github.com/adityak368/swissknife/middleware/tracing"
        "github.com/adityak368/swissknife/middleware/ratelimiter"
        "github.com/adityak368/swissknife/localization/i18n/middleware"
        "github.com/adityak368/swissknife/middleware/errorhandler"
        "github.com/labstack/echo/v4"
    )

    // Using Echo Middleware
    e := echo.New()

    e.Use(tracing.EchoTracingMiddleware(tracer, "AppName"))
    e.Use(middleware.EchoLocalizer())
    e.Use(ratelimiter.RateLimitMiddleware())
    e.HTTPErrorHandler = errorhandler.EchoHTTPErrorHandlerMiddleware
```

### ObjectStore

- Supplies Helpers to Upload/Download File/Image to Amazon S3 Store

```go
    import (
        "github.com/adityak368/swissknife/objectstore"
        "github.com/adityak368/swissknife/objectstore/s3"
        "github.com/labstack/echo/v4"
    )

    store := s3.New()

    func UploadAvatar(c echo.Context) error {
       // Using echo context c
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

        bucket := "abcd"
        id := "id"

        uploadResult, err := store.AddImage(bucket, id, bufReader)
        if err != nil {
            return err
        }

	    return c.JSON(http.StatusOK, "OK")
    }
```

### Response

- Defines the input/output interfaces for the internal business handlers
- Each internal function would return a result and error

```go
    import "github.com/adityak368/swissknife/response"

    func Login(c context.Context, user *User) (response.Result, error) {
        if user.Authenticate() {
            return &response.ExecResult{Result: map[string]string{token: "Token"}, MessageID: "UserAuthenticated"}, nil
        }
        return nil, response.NewError(http.StatusBadRequest, "InvalidEmailPassword")
    }
```

### Validation

- Provides helpers to perform input validation
- Uses go-playground for validation

```go
    import "github.com/adityak368/swissknife/validation/playground"

    type Test struct {
        A int `validate:"required,min=8,max=16"`
    }
    t := Test{A:10}
    validator := playground.New()
    validator.Validate(t)
```
