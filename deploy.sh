go build
ssh simplecloud "sudo stop gamedev"
rsync -avzP --progress game simplecloud:/home/god/sites/game.dev.zhuharev.ru/app/game
ssh simplecloud "sudo start gamedev"
