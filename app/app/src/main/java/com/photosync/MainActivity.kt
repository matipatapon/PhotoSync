package com.photosync

import android.content.Intent
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.text.input.TextFieldLineLimits
import androidx.compose.foundation.text.input.rememberTextFieldState
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
import androidx.compose.ui.unit.dp
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
                Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
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
                .padding(50.dp)
                .fillMaxSize(),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.spacedBy(10.dp, Alignment.CenterVertically),
            content = {
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
            Text(folder)
        }
        Text(info)
        Button(
            content= {Text("Add folder")},
            onClick = { addFolderToSync() }
        )
        Button(
            content={Text("Sync")},
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
