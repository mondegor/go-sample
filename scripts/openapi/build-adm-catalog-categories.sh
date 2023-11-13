
mrcmd_func_openapi_build_adm_catalog_categories() {
  local sectionDir="${1:?}" # sample: .../admin-api, .../public-api
  local sharedDir="${2:?}" # sample: .../_shared
  local unitDir="${sectionDir}/catalog/categories"

  OPENAPI_HEAD_PATH="${unitDir}/head.yaml"

  OPENAPI_SERVERS=(
    "${sectionDir}/servers.yaml"
  )

  OPENAPI_TAGS=(
    "${unitDir}/tags.yaml"
  )

  OPENAPI_PATHS=(
    "${unitDir}/paths.yaml"
  )

#  OPENAPI_COMPONENTS_HEADERS=(
#    "${sharedDir}/components/headers/"
#  )

  OPENAPI_COMPONENTS_PARAMETERS=(
    "${sharedDir}/components/parameters/App.Request.Header.AcceptLanguage.yaml"
    "${sharedDir}/components/parameters/App.Request.Header.CorrelationID.yaml"
    "${sharedDir}/components/parameters/App.Request.Query.Filter.Statuses.yaml"
    "${sharedDir}/components/parameters/App.Request.Query.Filter.SearchText.yaml"
    "${sharedDir}/components/parameters/App.Request.Query.ListPager.yaml"
    "${sharedDir}/components/parameters/App.Request.Query.ListSorter.yaml"

    "${unitDir}/components-parameters.yaml"
  )

  OPENAPI_COMPONENTS_SCHEMAS=(
    # "${sharedDir}/components/schemas/fields/App.Field.Article.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.Caption.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.CreatedAt.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.ImageURL.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.IntegerID.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.ListPager.Total.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.Price.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.Status.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.StringID.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.UpdatedAt.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.UUID.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.Version.yaml"
    "${sharedDir}/components/schemas/App.Request.Model.ChangeStatus.yaml"
    # "${sharedDir}/components/schemas/App.Request.Model.MoveItem.yaml"
    "${sharedDir}/components/schemas/App.Response.Model.BinaryFile.yaml"
    "${sharedDir}/components/schemas/App.Response.Model.CreateItem.yaml"
    "${sharedDir}/components/schemas/App.Response.Model.Error.yaml"
    "${sharedDir}/components/schemas/App.Response.Model.FileInfo.yaml"

    "${unitDir}/components-schemas.yaml"
  )

  OPENAPI_COMPONENTS_RESPONSES=(
    "${sharedDir}/components/responses/App.Response.Errors.yaml"
    "${sharedDir}/components/responses/App.Response.ErrorsAuth.yaml"

    "${unitDir}/components-responses.yaml"
  )

  OPENAPI_SECURITY_SCHEMES=(
    "${sharedDir}/securitySchemes.yaml"
  )
}
