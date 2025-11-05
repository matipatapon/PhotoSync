package com.photosync.ui.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color

private val LightColorScheme = lightColorScheme(
    primary = White,
    onPrimary = Black,
    secondary = Black,
    onSecondary = White,
    background = BackgroundColor,
    tertiary = Color.Cyan,
    error = Color(0xFF460000),
    onError = Color(0xFFFF0000),
    errorContainer = Color(0xFF330000)
)

@Composable
fun AppTheme(
    content: @Composable () -> Unit
) {
    MaterialTheme(
        colorScheme = LightColorScheme,
        typography = Typography,
        content = content
    )
}
