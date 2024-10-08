
function mrcmd_plugins_openapi_method_init() {
  readonly OPENAPI_CAPTION="Go Sample REST API Builder"

  readonly OPENAPI_VARS=(
    "OPENAPI_SRC_DIR"
    "OPENAPI_SRC_SHARED_DIR"
    "OPENAPI_ASSEMBLY_DIR"
    "OPENAPI_FILENAME_PREFIX"
  )

  # default values of array: OPENAPI_VARS
  readonly OPENAPI_VARS_DEFAULT=(
    "${APPX_DIR}/docs/api-src"
    "${APPX_DIR}/docs/api-src/_shared"
    "${APPX_DIR}/docs/api"
    "openapi-"
  )

  mrcore_dotenv_init_var_array OPENAPI_VARS[@] OPENAPI_VARS_DEFAULT[@]
}

function mrcmd_plugins_openapi_method_config() {
  mrcore_dotenv_echo_var_array OPENAPI_VARS[@]
}

function mrcmd_plugins_openapi_method_export_config() {
  mrcore_dotenv_export_var_array OPENAPI_VARS[@]
}

function mrcmd_plugins_openapi_method_exec() {
  local currentCommand="${1:?}" # script name: adm.sh, public.sh, ...
  local sectionName="" # sample: admin-api, public-api, ...
  local fileNamePostfix="" # sample: all, main, ...

  case "${currentCommand}" in

    build-all)
      mrcmd openapi build-adm-all
      mrcmd openapi build-int-all
      mrcmd openapi build-pub-all
      ;;

    build-full)
      mrcmd openapi build-adm
      mrcmd openapi build-int
      mrcmd openapi build-pub
      ;;

    build-adm-all)
      mrcmd openapi build-adm
      mrcmd openapi build-adm-catalog
      ;;

    build-int-all)
      mrcmd openapi build-int
      ;;

    build-pub-all)
      mrcmd openapi build-pub
      mrcmd openapi build-pub-catalog
      mrcmd openapi build-pub-file-station
      ;;

    build-adm)
      sectionName="admin-api"
      ;;

    build-int)
      sectionName="internal-api"
      ;;

    build-adm-catalog)
      sectionName="admin-api"
      fileNamePostfix="catalog"
      ;;

    build-pub)
      sectionName="public-api"
      ;;

    build-pub-catalog)
      sectionName="public-api"
      fileNamePostfix="catalog"
      ;;

    build-pub-file-station)
      sectionName="public-api"
      fileNamePostfix="file-station"
      ;;

    *)
      ${RETURN_FALSE}
      ;;

  esac

  if [ -n "${sectionName}" ]; then
    mrcore_import "${MRCMD_PLUGINS_DIR}/libs/openapi-lib.sh"
    openapi_lib_build_apidoc \
      "openapi/${currentCommand}" \
      "${OPENAPI_SRC_DIR}" \
      "${OPENAPI_SRC_SHARED_DIR}" \
      "${OPENAPI_ASSEMBLY_DIR}" \
      "${sectionName}" \
      "${OPENAPI_FILENAME_PREFIX}" \
      "${fileNamePostfix}"
  fi
}

function mrcmd_plugins_openapi_method_help() {
  #markup:"|-|-|---------|-------|-------|---------------------------------------|"
  echo -e "${CC_YELLOW}Commands:${CC_END}"
  echo -e "    build-all                 Builds all API docs"
  echo -e "    build-full                Builds only full API docs"
  echo -e "    build-adm-all             Builds all admin API docs"
  echo -e "    build-adm                 Builds only full admin API docs"
  echo -e "    build-adm-catalog         Builds admin Catalog API docs"
  echo -e "    build-int-all             Builds internal (system) API docs"
  echo -e "    build-pub-all             Builds all public API docs"
  echo -e "    build-pub                 Builds only full public API docs"
  echo -e "    build-pub-catalog         Builds public Catalog API docs"
}
