brew install android-commandlinetools

sdkmanager "platform-tools" "platforms;android-28"

export $ANDROID_HOME=/usr/local/share/android-commandlinetools
on my linux box it's in my home dir under ~/Android and some people on Mac have it under ~/Library/Android

brew install android-ndk

go install github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile@latest

ebitenmobile bind -target android -javapkg org.sinisterstuf.cr1cktbin -o cr1ckt.aar ./mobile/

more info: https://ebiten.org/documents/mobile.html#Android

also had to change ldtkgo Ebiten renderer to never use ebitenutil.ImageFromFile because in mobile there's not filesystem access and this function doesn't exist

to build the apk from the Android projet run ./gradlew assemble
