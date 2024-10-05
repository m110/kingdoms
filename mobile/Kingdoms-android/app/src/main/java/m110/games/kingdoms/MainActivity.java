package m110.games.kingdoms;

import androidx.appcompat.app.AppCompatActivity;

import android.os.Bundle;
import go.Seq;
import m110.kingdoms.mobile.EbitenView;

public class MainActivity extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        Seq.setContext(getApplicationContext());

    }

    private EbitenView getEbitenView() {
        return (EbitenView)this.findViewById(R.id.ebitenView);
    }

    @Override
    protected void onPause() {
        super.onPause();
        this.getEbitenView().suspendGame();
    }

    @Override
    protected void onResume() {
        super.onResume();
        this.getEbitenView().resumeGame();
    }
}