go build
ssh simplecloud "sudo stop gamedev"
scp game simplecloud:/home/god/sites/game.dev.zhuharev.ru/app/
ssh simplecloud "sudo start gamedev"