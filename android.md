brew install android-commandlinetools

sdkmanager "platform-tools" "platforms;android-28"

export $ANDROID_HOME=/usr/local/share/android-commandlinetools

brew install android-ndk

ebitenmobile bind -target android -javapkg org.sinisterstuf.cr1ckt -o cr1ckt.aar ./mobile/

more info: https://ebiten.org/documents/mobile.html#Android

also had to change ldtkgo Ebiten renderer to never use ebitenutil.ImageFromFile because in mobile there's not filesystem access and this function doesn't exist
