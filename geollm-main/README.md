# README

geollm后端部分

## 项目结构

```
.
├── Dockerfile
├── README.md
├── api
│   ├── answersheet.go
│   ├── class.go
│   ├── exam.go
│   ├── student.go
│   ├── user.go
│   └── xueqing.go
├── conf
│   └── config.yaml
├── dao
│   ├── answersheetDao.go
│   ├── classDao.go
│   ├── examDao.go
│   ├── studentDao.go
│   ├── userDao.go
│   └── xueqing.go
├── go.mod
├── go.sum
├── main.go
├── middleware
│   ├── cors.go
│   └── jwt.go
├── model
│   ├── answersheet.go
│   ├── class.go
│   ├── db.go
│   ├── exam.go
│   ├── student.go
│   ├── user.go
│   └── xueqing.go
├── routers
│   └── router.go
├── tmp
│   ├── runner-build
│   └── runner-build-errors.log
└── utils
    ├── GeneralReturn.go
    ├── cfg.go
    ├── docx.go
    ├── errmsg.go
    ├── evaluator.go
    ├── minio.go
    ├── ocr.go
    ├── token.go
    └── validator.go
```

## 接口文档

[Apifox](https://apifox.com/apidoc/shared-7fcf9b93-3220-4290-9cbf-04d42566304a)

## ToDo
- [ ] 日志导出
- [ ] 更改Token为Session