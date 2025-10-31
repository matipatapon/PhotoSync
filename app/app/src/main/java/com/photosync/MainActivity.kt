package com.photosync

import android.content.Intent
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxHeight
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.sizeIn
import androidx.compose.foundation.layout.wrapContentHeight
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.text.input.TextFieldLineLimits
import androidx.compose.foundation.text.input.rememberTextFieldState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.Button
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SecureTextField
import androidx.compose.material3.Text
import androidx.compose.material3.TextField
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Modifier
import com.photosync.ui.theme.AppTheme
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.room.Room
import com.photosync.api.ApiHandler
import com.photosync.database.LocalDatabase
import com.photosync.view_models.FolderViewModel
import com.photosync.view_models.LoginViewModel
import com.photosync.view_models.Window

class MainActivity : ComponentActivity() {
    private var localDatabase: LocalDatabase? = null
    private var loginViewModel: LoginViewModel? = null
    private var folderViewModel: FolderViewModel? = null
    private val apiHandler: ApiHandler = ApiHandler()

    override fun onCreate(savedInstanceState: Bundle?) {
        localDatabase = Room.databaseBuilder(
            applicationContext,
            LocalDatabase::class.java, "PhotoSync"
        ).build()
        loginViewModel = LoginViewModel(localDatabase!!, apiHandler)
        folderViewModel = FolderViewModel(localDatabase!!, application, apiHandler)

        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            AppTheme {
                Scaffold(modifier = Modifier.fillMaxSize(),
                    ) { innerPadding ->
                    View(
                        innerPadding
                    )
                }
            }
        }
    }

    @Composable
    fun View(innerPadding: PaddingValues){
        val window by loginViewModel!!.window.collectAsState()
        Column(
            modifier = Modifier
                .fillMaxSize().padding(50.dp).verticalScroll(rememberScrollState()),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.spacedBy(10.dp, Alignment.CenterVertically),
            content = {
                Text(text = "PhotoSync", Modifier.fillMaxWidth().wrapContentHeight(
                    Alignment.Bottom), textAlign = TextAlign.Center, fontSize = 24.sp, fontWeight = FontWeight.Bold)
                Spacer(Modifier.weight(0.5f))
                when (window) {
                    Window.Load -> {
                        loginViewModel!!.load()
                    }
                    Window.Login -> {
                        LoginForm()
                    }
                    Window.Sync -> {
                        Folders()
                    }
                }
                Spacer(Modifier.weight(0.5f))
            }
        )
    }

    fun addFolderToSync() {
        val intent = Intent(Intent.ACTION_OPEN_DOCUMENT_TREE)
        addFolderLauncher.launch(intent)
    }

    var addFolderLauncher = registerForActivityResult(ActivityResultContracts.StartActivityForResult()) { result ->
        if (result.resultCode == RESULT_OK) {
             result.data?.data?.let {
                uri ->
                    val contentResolver = applicationContext.contentResolver
                    val takeFlags: Int = Intent.FLAG_GRANT_READ_URI_PERMISSION
                    contentResolver.takePersistableUriPermission(uri, takeFlags)
                    folderViewModel!!.addFolderToSync(uri)
            }
        }
    }

    @Composable
    fun Folders(){
        val folders by folderViewModel!!.folders.collectAsState()
        val info by folderViewModel!!.info.collectAsState()
        for(folder in folders){
            Row(modifier = Modifier.fillMaxWidth(),
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.SpaceEvenly,
                content={
                    Text(text=folder, textAlign = TextAlign.Center)
                    Spacer(modifier = Modifier.weight(1f))
                    Button(modifier = Modifier.fillMaxHeight(), content={ Text("X") }, onClick = {})
                }
            )
        }
        Button(
            content= {Text("+")},
            modifier = Modifier.fillMaxWidth(),
            onClick = { addFolderToSync() }
        )
        Spacer(Modifier.height(25.dp))
        if(info != ""){
            Text(text=info)
        }
        Button(
            content={Text("Sync")},
            modifier = Modifier.fillMaxWidth(),
            onClick = {
                folderViewModel!!.syncFolders()
            }
        )
    }

    @Composable
    fun LoginForm() {
        val appSettings = loginViewModel!!.appSettings
        var initialServer = ""
        var initialUsername = ""
        if(appSettings != null){
            initialServer = appSettings.server
            initialUsername = appSettings.login
        }
        val server = rememberTextFieldState(initialText = initialServer)
        val username = rememberTextFieldState(initialText = initialUsername)
        val password = rememberTextFieldState(initialText = "")
        val loginStatus by loginViewModel!!.loginStatus.collectAsState()
        TextField(
            state = server,
            placeholder = { Text("server") },
            lineLimits = TextFieldLineLimits.SingleLine,
            modifier = Modifier.fillMaxWidth(),
            enabled = !loginStatus.isPending()
        )
        TextField(
            state = username,
            placeholder = { Text("login") },
            lineLimits = TextFieldLineLimits.SingleLine,
            modifier = Modifier.fillMaxWidth(),
            enabled = !loginStatus.isPending()
        )
        SecureTextField(
            state = password,
            placeholder = { Text("password") },
            modifier = Modifier.fillMaxWidth(),
            enabled = !loginStatus.isPending(),
        )
        Text(loginStatus.getError())
        Button(
            onClick = {
                loginViewModel!!.login(server.text.toString(), username.text.toString(), password.text.toString())
            },
            enabled = !loginStatus.isPending(),
            content = {
                Text("Login")
            },
            modifier = Modifier.fillMaxWidth(),
        )
    }
}
