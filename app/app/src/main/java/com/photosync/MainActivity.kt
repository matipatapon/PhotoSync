package com.photosync

import android.content.Intent
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.sizeIn
import androidx.compose.foundation.layout.wrapContentHeight
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.input.TextFieldLineLimits
import androidx.compose.foundation.text.input.TextFieldState
import androidx.compose.foundation.text.input.rememberTextFieldState
import androidx.compose.foundation.text.selection.LocalTextSelectionColors
import androidx.compose.foundation.text.selection.TextSelectionColors
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SecureTextField
import androidx.compose.material3.Text
import androidx.compose.material3.TextField
import androidx.compose.material3.TextFieldDefaults
import androidx.compose.runtime.Composable
import androidx.compose.runtime.CompositionLocalProvider
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Modifier
import com.photosync.ui.theme.AppTheme
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.core.net.toUri
import androidx.room.Room
import com.photosync.api.ApiHandler
import com.photosync.database.LocalDatabase
import com.photosync.ui.theme.Black
import com.photosync.ui.theme.DisabledContainerColor
import com.photosync.view_models.FolderStatus
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
    fun Header(){
        Text(text = "PhotoSync",
            color = MaterialTheme.colorScheme.primary,
            modifier = Modifier
            .fillMaxWidth()
            .wrapContentHeight(
                Alignment.Bottom
            ), textAlign = TextAlign.Center, fontSize = 24.sp, fontWeight = FontWeight.Bold)
    }
    
    @Composable
    fun View(innerPadding: PaddingValues){
        val window by loginViewModel!!.window.collectAsState()
        val customTextSelectionColors = TextSelectionColors(
            handleColor = Black,
            backgroundColor = Black.copy(alpha = 0.4f)
        )
        CompositionLocalProvider(LocalTextSelectionColors provides customTextSelectionColors) {
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
    fun MyTextFiled(text: String, enabled: Boolean, state: TextFieldState, secure: Boolean = false){
        val colors = TextFieldDefaults.colors(
            unfocusedContainerColor = MaterialTheme.colorScheme.primary,
            focusedContainerColor = MaterialTheme.colorScheme.primary ,
            unfocusedTextColor = MaterialTheme.colorScheme.onPrimary,
            focusedTextColor = MaterialTheme.colorScheme.onPrimary,
            cursorColor = MaterialTheme.colorScheme.onPrimary,
            focusedIndicatorColor = MaterialTheme.colorScheme.tertiary,
            disabledContainerColor = DisabledContainerColor
        )
        if(secure){
            SecureTextField(
                state = state,
                placeholder = { Text(text) },
                modifier = Modifier.fillMaxWidth(),
                enabled = enabled,
                colors = colors
            )
        } else{
            TextField(
                state = state,
                placeholder = { Text(text) },
                lineLimits = TextFieldLineLimits.SingleLine,
                modifier = Modifier.fillMaxWidth(),
                enabled = enabled,
                colors = colors
            )
        }
    }

    @Composable
    fun MyButton(text: String, enabled: Boolean, onClick: ()-> Unit){
        Button(
            onClick = onClick,
            enabled = enabled,
            content = {
                Text(text)
            },
            modifier = Modifier.fillMaxWidth(),
            colors = ButtonDefaults.buttonColors(
                contentColor = MaterialTheme.colorScheme.onPrimary,
                containerColor = MaterialTheme.colorScheme.primary
            )
        )
    }

    @Composable
    private fun Popup(){
        val folderStatus = folderViewModel!!.status.collectAsState()
        if(folderStatus.value.type == FolderStatus.Type.Idle){
            return
        }
        Column(
            modifier = Modifier
                .fillMaxSize()
                .background(Color(0, 0, 0, 203)),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Center,
            content = {
                val columnColor = if (folderStatus.value.type == FolderStatus.Type.Error) MaterialTheme.colorScheme.error else MaterialTheme.colorScheme.primary
                Column(
                    modifier = Modifier
                        .clip(RoundedCornerShape(10.dp))
                        .sizeIn(maxWidth = 300.dp, maxHeight = 400.dp)
                        .background(color = columnColor)
                        .padding(10.dp),
                    horizontalAlignment = Alignment.CenterHorizontally,
                    verticalArrangement = Arrangement.spacedBy(10.dp, Alignment.CenterVertically),
                    ){
                    if(folderStatus.value.type == FolderStatus.Type.Sync){
                        Text(text = "Uploading",
                            color = MaterialTheme.colorScheme.onPrimary,
                            textAlign = TextAlign.Center,
                            fontSize = 20.sp)
                        if(folderStatus.value.info != ""){
                            Text(text = folderStatus.value.info,
                                color = MaterialTheme.colorScheme.onPrimary,
                                textAlign = TextAlign.Center,
                                maxLines = 1)
                        }
                    } else if(folderStatus.value.type == FolderStatus.Type.Error){
                        Text(text = "Error",
                            color = MaterialTheme.colorScheme.onError,
                            textAlign = TextAlign.Center,
                            fontSize = 20.sp)
                        Text(text = folderStatus.value.info,
                            color = MaterialTheme.colorScheme.onError,
                            textAlign = TextAlign.Center,
                            maxLines = 1)
                        Button(
                            content = {Text("Ok")},
                            onClick = {
                                folderViewModel!!.resetStatus()
                            },
                            colors = ButtonDefaults.buttonColors(
                                contentColor = MaterialTheme.colorScheme.onError,
                                containerColor = MaterialTheme.colorScheme.errorContainer
                            )
                        )
                    }
                }
            }
        )
    }

    @Composable
    fun Folders(){
        val folders by folderViewModel!!.folders.collectAsState()
        val folderStatus = folderViewModel!!.status.collectAsState()
        val enabled = folderStatus.value.type == FolderStatus.Type.Idle
        Box(content= {
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(50.dp)
                    .verticalScroll(rememberScrollState()),
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(10.dp, Alignment.CenterVertically),
                content = {
                    Header()
                    Spacer(Modifier.weight(0.5f))
                    for (folder in folders) {
                        Row(
                            modifier = Modifier
                                .fillMaxWidth()
                                .clip(RoundedCornerShape(25.dp))
                                .background(MaterialTheme.colorScheme.primary)
                                .padding(10.dp),
                            verticalAlignment = Alignment.CenterVertically,
                            content = {
                                Spacer(Modifier.weight(0.5f))
                                Text(
                                    text = folder.uri.toUri().path.toString().substringAfter(":"),
                                    textAlign = TextAlign.Center,
                                    maxLines = 1
                                )
                                Spacer(Modifier.weight(0.5f))
                                Button(
                                    content = { Text("x") },
                                    onClick = {
                                        folderViewModel!!.deleteFolder(folder)
                                    },
                                    colors = ButtonDefaults.buttonColors(
                                        contentColor = MaterialTheme.colorScheme.onSecondary,
                                        containerColor = MaterialTheme.colorScheme.secondary
                                    ),
                                    enabled = enabled
                                )
                            }
                        )
                    }
                    Spacer(Modifier.weight(0.5f))
                    MyButton(
                        text = "+",
                        enabled = enabled,
                        onClick = { addFolderToSync() }
                    )
                    MyButton(
                        text = "Sync",
                        enabled = enabled,
                        onClick = { folderViewModel!!.syncFolders()}
                    )
                }
            )
            Popup()
        })
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
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(50.dp)
                .verticalScroll(rememberScrollState()),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.spacedBy(10.dp, Alignment.CenterVertically),
            content = {
                Header()
                Spacer(Modifier.weight(0.5f))
                MyTextFiled(text = "server", enabled = !loginStatus.isPending(), state = server)
                MyTextFiled(text = "login", enabled = !loginStatus.isPending(), state = username)
                MyTextFiled(text = "password", enabled = !loginStatus.isPending(), state = password, secure = true)
                val error = loginStatus.getError()
                if(error != ""){
                    Text(
                        text=error,
                        color = MaterialTheme.colorScheme.onError,
                        textAlign = TextAlign.Center,
                        modifier = Modifier
                            .background(MaterialTheme.colorScheme.error)
                            .fillMaxWidth()
                            .padding(10.dp)
                    )
                }
                MyButton(text = "Login", !loginStatus.isPending(), onClick = {
                    loginViewModel!!.login(
                        server.text.toString(),
                        username.text.toString(),
                        password.text.toString()
                    )
                })
                Spacer(Modifier.weight(0.5f))
            }
        )
    }
}
