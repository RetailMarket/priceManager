if [[ ! -e out/ ]];
		then
			mkdir out/
		fi 

		go build -o out/build app/priceManager/main.go; ./out/build
