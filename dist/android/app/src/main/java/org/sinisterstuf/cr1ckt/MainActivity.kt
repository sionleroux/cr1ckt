package org.sinisterstuf.cr1ckt

import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import androidx.core.view.WindowCompat.getInsetsController
import androidx.core.view.WindowInsetsCompat
import androidx.core.view.WindowInsetsControllerCompat.BEHAVIOR_SHOW_TRANSIENT_BARS_BY_SWIPE
import go.Seq
import org.sinisterstuf.cr1cktbin.mobile.EbitenView

class MainActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        Seq.setContext(applicationContext)

        // Immersive full-screen mode
        val controller = getInsetsController(window, window.decorView)
        controller.systemBarsBehavior = BEHAVIOR_SHOW_TRANSIENT_BARS_BY_SWIPE
        controller.hide(WindowInsetsCompat.Type.systemBars())

        setContentView(R.layout.activity_main)
    }

    override fun onPause() {
        findViewById<EbitenView>(R.id.ebiten).suspendGame()
        super.onPause()
    }

    override fun onResume() {
        findViewById<EbitenView>(R.id.ebiten).resumeGame()
        super.onResume()
    }
}
