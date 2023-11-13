
mrcmd_func_openapi_build_adm_catalog() {
  local sectionDir="${1:?}" # sample: .../admin-api, .../public-api
  local sharedDir="${2:?}" # sample: .../_shared
  local moduleDir="${sectionDir}/catalog"

  local catalogCategoriesDir="${moduleDir}/categories"
  local catalogProductsDir="${moduleDir}/products"
  local catalogTrademarksDir="${moduleDir}/trademarks"

  OPENAPI_HEAD_PATH="${moduleDir}/head.yaml"

  OPENAPI_SERVERS=(
    "${sectionDir}/servers.yaml"
  )

  OPENAPI_TAGS=(
    "${catalogCategoriesDir}/tags.yaml"
    "${catalogProductsDir}/tags.yaml"
    "${catalogTrademarksDir}/tags.yaml"
  )

  OPENAPI_PATHS=(
    "${catalogCategoriesDir}/paths.yaml"
    "${catalogProductsDir}/paths.yaml"
    "${catalogTrademarksDir}/paths.yaml"
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

    "${catalogCategoriesDir}/components-parameters.yaml"
    "${catalogProductsDir}/components-parameters.yaml"
    "${catalogTrademarksDir}/components-parameters.yaml"
  )

  OPENAPI_COMPONENTS_SCHEMAS=(
    "${sharedDir}/components/schemas/fields/App.Field.Article.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.Caption.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.CreatedAt.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.ImageURL.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.IntegerID.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.ListPager.Total.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.Price.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.Status.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.StringID.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.UpdatedAt.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.UUID.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.Version.yaml"
    "${sharedDir}/components/schemas/App.Request.Model.ChangeStatus.yaml"
    "${sharedDir}/components/schemas/App.Request.Model.MoveItem.yaml"
    "${sharedDir}/components/schemas/App.Response.Model.BinaryFile.yaml"
    "${sharedDir}/components/schemas/App.Response.Model.CreateItem.yaml"
    "${sharedDir}/components/schemas/App.Response.Model.Error.yaml"
    "${sharedDir}/components/schemas/App.Response.Model.FileInfo.yaml"

    "${catalogCategoriesDir}/components-schemas.yaml"
    "${catalogProductsDir}/components-schemas.yaml"
    "${catalogTrademarksDir}/components-schemas.yaml"
  )

  OPENAPI_COMPONENTS_RESPONSES=(
    "${sharedDir}/components/responses/App.Response.Errors.yaml"
    "${sharedDir}/components/responses/App.Response.ErrorsAuth.yaml"

    "${catalogCategoriesDir}/components-responses.yaml"
    "${catalogProductsDir}/components-responses.yaml"
    "${catalogTrademarksDir}/components-responses.yaml"
  )

  OPENAPI_SECURITY_SCHEMES=(
    "${sharedDir}/securitySchemes.yaml"
  )
}
