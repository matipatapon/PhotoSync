package com.photosync

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.text.input.TextFieldLineLimits
import androidx.compose.foundation.text.input.rememberTextFieldState
import androidx.compose.foundation.text.input.setTextAndPlaceCursorAtEnd
import androidx.compose.material3.Button
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SecureTextField
import androidx.compose.material3.Text
import androidx.compose.material3.TextField
import androidx.compose.runtime.Composable
import androidx.compose.runtime.mutableStateOf
import androidx.compose.ui.Modifier
import com.photosync.ui.theme.AppTheme
import androidx.compose.runtime.remember
import androidx.compose.runtime.getValue
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.unit.dp

class MainActivity : ComponentActivity() {
    private val serverRepository: ServerRepository = ServerRepository()
    private val stage: Int = 1;

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            AppTheme {
                Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
                    LoginForm(
                        innerPadding
                    )
                }
            }
        }
    }

    @Composable
    fun LoginForm(innerPadding: PaddingValues) {
        val username = rememberTextFieldState(initialText = "")
        val password = rememberTextFieldState(initialText = "")
        var showButton by remember {mutableStateOf(true)}
        Column(
            Modifier
                .fillMaxSize()
                .padding(innerPadding),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.spacedBy(10.dp, Alignment.CenterVertically),
            content = {
                TextField(
                    state = username,
                    placeholder = { Text("login") },
                    lineLimits = TextFieldLineLimits.SingleLine,
                    modifier = Modifier
                )
                SecureTextField(
                    state = password,
                    placeholder = { Text("password") },
                    modifier = Modifier,
                )
                Button(
                    onClick = {
                        showButton = false
                        password.setTextAndPlaceCursorAtEnd("haha")
                    },
                    enabled = showButton,
                    content = {
                        Text("Login")
                    }
                )
            })
    }
}
