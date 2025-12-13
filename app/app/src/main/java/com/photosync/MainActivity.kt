package com.photosync

import android.content.Intent
import android.os.Bundle
import android.widget.Space
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.Image
import androidx.compose.foundation.background
import androidx.compose.foundation.horizontalScroll
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.sizeIn
import androidx.compose.foundation.layout.width
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
import androidx.compose.ui.graphics.RectangleShape
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.core.net.toUri
import androidx.room.Room
import com.photosync.api.ApiHandler
import com.photosync.database.LocalDatabase
import com.photosync.ui.theme.Background
import com.photosync.ui.theme.Black
import com.photosync.ui.theme.DisabledContainerColor
import com.photosync.view_models.FolderStatus
import com.photosync.view_models.FolderViewModel
import com.photosync.view_models.LoginViewModel
import com.photosync.view_models.Window
import com.example.app.R
import com.photosync.ui.theme.LightGray

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
            .fillMaxWidth(),
            textAlign = TextAlign.Center, fontSize = 24.sp, fontWeight = FontWeight.Bold)
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
            unfocusedContainerColor = MaterialTheme.colorScheme.secondary,
            focusedContainerColor = MaterialTheme.colorScheme.secondary ,
            unfocusedTextColor = MaterialTheme.colorScheme.onSecondary,
            focusedTextColor = MaterialTheme.colorScheme.onSecondary,
            cursorColor = MaterialTheme.colorScheme.onSecondary,
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
                .background(Background)
                .padding(50.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Center,
            content = {
                Header()
                Spacer(Modifier.height(20.dp))
                Spacer(modifier = Modifier.weight(0.3f))
                Column(
                    modifier = Modifier
                        .clip(RoundedCornerShape(10.dp))
                        .background(color = LightGray)
                        .padding(10.dp)
                        .fillMaxWidth(),
                    horizontalAlignment = Alignment.CenterHorizontally,
                    verticalArrangement = Arrangement.spacedBy(10.dp, Alignment.CenterVertically),
                    ){
                    if(folderStatus.value.type == FolderStatus.Type.Sync){
                        Text(text = "Uploading",
                            color = MaterialTheme.colorScheme.primary,
                            textAlign = TextAlign.Center,
                            fontSize = 20.sp)
                        if(folderStatus.value.info != ""){
                            Text(text = folderStatus.value.info,
                                color = MaterialTheme.colorScheme.primary,
                                textAlign = TextAlign.Center,
                                fontSize = 12.sp)
                        }
                    } else if(folderStatus.value.type == FolderStatus.Type.Error){
                        Text(text = "Error",
                            color = MaterialTheme.colorScheme.primary,
                            textAlign = TextAlign.Center,
                            fontSize = 20.sp)
                        Text(text = folderStatus.value.info,
                            color = MaterialTheme.colorScheme.primary,
                            textAlign = TextAlign.Left,
                            fontSize = 12.sp)
                        MyButton(
                            text = "Ok",
                            enabled = true,
                            onClick = { folderViewModel!!.resetStatus() }
                        )
                    } else if(folderStatus.value.type == FolderStatus.Type.Confirmation) {
                        Text(text = "Finished",
                            color = MaterialTheme.colorScheme.primary,
                            textAlign = TextAlign.Center,
                            fontSize = 20.sp)
                        Spacer(modifier = Modifier.height(10.dp))
                        MyButton(
                            text = "Ok",
                            enabled = true,
                            onClick = { folderViewModel!!.resetStatus() }
                        )
                    }
                }
                Spacer(modifier = Modifier.weight(0.7f))
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
                    .padding(20.dp),
                content = {
                    Spacer(Modifier.height(30.dp))
                    Header()
                    Spacer(Modifier.height(20.dp))
                    Column(
                        modifier = Modifier
                            .verticalScroll(rememberScrollState())
                            .clip(RoundedCornerShape(10.dp))
                            .background(LightGray)
                            .weight(1f, true)
                            .padding(20.dp),
                        horizontalAlignment = Alignment.CenterHorizontally,
                        verticalArrangement = Arrangement.spacedBy(10.dp),
                        content= {
                            for (folder in folders) {
                                Row(
                                    modifier = Modifier
                                        .fillMaxWidth()
                                        .clip(RoundedCornerShape(5.dp))
                                        .background(Background),
                                    verticalAlignment = Alignment.CenterVertically,
                                    content = {
                                        Spacer(Modifier.width(20.dp))
                                        Text(
                                            text = folder.uri.toUri().path.toString()
                                                .substringAfter(":"),
                                            maxLines = 1,
                                            color = MaterialTheme.colorScheme.primary,
                                            fontSize = 16.sp,
                                            modifier = Modifier.horizontalScroll(rememberScrollState()).weight(1f, fill = true),
                                        )
                                        Spacer(Modifier.width(20.dp))
                                        Button(
                                            modifier = Modifier.width(40.dp).height(40.dp),
                                            contentPadding = PaddingValues(5.dp),
                                            content = {
                                                Image(
                                                    painter = painterResource(id = R.drawable.trash),
                                                    contentDescription = null
                                                )
                                            },
                                            onClick = {
                                                folderViewModel!!.deleteFolder(folder)
                                            },
                                            colors = ButtonDefaults.buttonColors(
                                                contentColor = MaterialTheme.colorScheme.onPrimary,
                                            ),
                                            enabled = enabled,
                                            shape = RectangleShape,
                                        )
                                    }
                                )
                            }
                    })
                    Spacer(Modifier.height(10.dp))
                    Text(text = "256 files are not synchronized",
                        color = MaterialTheme.colorScheme.primary,
                        modifier = Modifier
                            .fillMaxWidth(),
                        textAlign = TextAlign.Center, fontSize = 16.sp, fontWeight = FontWeight.Bold)
                    Spacer(Modifier.height(10.dp))
                    MyButton(
                        text = "Add folder",
                        enabled = enabled,
                        onClick = { addFolderToSync() }
                    )
                    Spacer(Modifier.height(5.dp))
                    MyButton(
                        text = "Synchronize",
                        enabled = enabled,
                        onClick = { folderViewModel!!.syncFolders()}
                    )
                    Spacer(Modifier.height(30.dp))
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
