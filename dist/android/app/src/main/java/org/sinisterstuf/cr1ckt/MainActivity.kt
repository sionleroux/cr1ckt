package org.sinisterstuf.cr1ckt

import android.os.Bundle
import android.view.Window
import android.view.WindowManager.LayoutParams.FLAG_FULLSCREEN
import androidx.appcompat.app.AppCompatActivity
import go.Seq
import org.sinisterstuf.cr1cktbin.mobile.EbitenView

class MainActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        Seq.setContext(applicationContext)
        requestWindowFeature(Window.FEATURE_NO_TITLE)
        window.setFlags(FLAG_FULLSCREEN, FLAG_FULLSCREEN)
        setContentView(R.layout.activity_main)
        supportActionBar?.hide()
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
