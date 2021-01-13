##
# SendReadable
#
# @file
# @version 0.1


css:
	cd assets;node_modules/postcss-cli/bin/postcss --use-autoprefixer -o dist/style.css src/css/style.css

assets: css
	pkger -o assets/compiled

build: assets
	go build
# end
