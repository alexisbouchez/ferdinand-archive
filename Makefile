# Run the application with watch mode
air:
	@go run github.com/cosmtrek/air@v1.51.0 -c ./zarf/air.toml

# Generate the CSS files
css:
	@npx tailwindcss -c ./zarf/tailwind.config.js -i ./views/app.css -o ./public/assets/styles.css --watch   

# Generate the templates
gentempl:
	@templ generate --watch --proxy="http://localhost:3000"

# Clean generated files
clean:
	@rm -rf ./views/**/**_templ.go ./views/**/**_templ.txt
