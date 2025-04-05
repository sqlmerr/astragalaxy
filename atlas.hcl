data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/model",
    "--dialect", "postgres", // | mysql | sqlite | sqlserver
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
