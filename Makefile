# variable to store the path of the frontend directory
FRONTEND_DIR = ~/Development/react/z-r-react-ui
FRONTEND_WEB_URL = http://localhost:5173

runbackdev:
	go run cmd/main.go dev

runfrontdev:
	sleep 2 && cd ${FRONTEND_DIR} && yarn dev

openweb:
	sleep 3 && open ${FRONTEND_WEB_URL}

# Tarea para ejecutar ambos proyectos en paralelo
rundev: 
	$(MAKE) -j 3 runbackdev runfrontdev openweb


# Configura la tarea por defecto
.DEFAULT_GOAL := rundev