package com.blackforestbytes.simplecloudnotifier2

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.offset
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Menu
import androidx.compose.material3.CenterAlignedTopAppBar
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FabPosition
import androidx.compose.material3.FloatingActionButton
import androidx.compose.material3.FloatingActionButtonDefaults
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.material3.TopAppBarDefaults
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.blackforestbytes.simplecloudnotifier2.ui.theme.Simplecloudnotifier2Theme

class MainActivity : ComponentActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            Simplecloudnotifier2Theme {
                // A surface container using the 'background' color from the theme
                Surface(modifier = Modifier.fillMaxSize(), color = MaterialTheme.colorScheme.background) {

                    Content()

                }
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun Content() {
    Scaffold(

        topBar = {
            CenterAlignedTopAppBar(
                colors = TopAppBarDefaults.centerAlignedTopAppBarColors(
                    containerColor = MaterialTheme.colorScheme.primaryContainer,
                    titleContentColor = MaterialTheme.colorScheme.primary,
                ),
                title = {
                    Text("Messages", maxLines = 1, overflow = TextOverflow.Ellipsis)
                },
                navigationIcon = {
                    IconButton(onClick = { /* do something */ }) {
                        Icon(painterResource(R.drawable.fas_gauge), contentDescription = "Menu", modifier = Modifier.size(24.dp))
                    }
                },

                actions = {
                    IconButton(onClick = { /* do something */ }) {
                        Icon(painterResource(R.drawable.fas_paper_plane_top), contentDescription = "Send message", modifier = Modifier.size(24.dp))
                    }
                },
            )
        },

        bottomBar = {
            NavigationBar {
                NavigationBarItem(
                    icon = { Icon(painterResource(R.drawable.fas_road), contentDescription = "Channels", modifier = Modifier.size(32.dp)) },
                    onClick = {},
                    selected = false,
                )
                NavigationBarItem(
                    icon = { Icon(painterResource(R.drawable.fas_computer), contentDescription = "Clients", modifier = Modifier.size(32.dp)) },
                    onClick = {},
                    selected = false,
                )
                NavigationBarItem(
                    icon = { Icon(painterResource(R.drawable.fas_key), contentDescription = "Keys", modifier = Modifier.size(32.dp)) },
                    onClick = {},
                    selected = false,
                )
                NavigationBarItem(
                    icon = { },
                    onClick = {},
                    selected = false,
                )
                NavigationBarItem(
                    icon = { Icon(painterResource(R.drawable.fas_bookmark), contentDescription = "Subscriptions", modifier = Modifier.size(32.dp)) },
                    onClick = {},
                    selected = false,
                )
                NavigationBarItem(
                    icon = { Icon(painterResource(R.drawable.fas_user), contentDescription = "User", modifier = Modifier.size(32.dp)) },
                    onClick = {},
                    selected = false,
                )
                NavigationBarItem(
                    icon = { Icon(painterResource(R.drawable.fas_gear), contentDescription = "Settings", modifier = Modifier.size(32.dp)) },
                    onClick = {},
                    selected = false,
                )
            }
        },

        floatingActionButton = {
            Box(){
                FloatingActionButton(
                    onClick = { /* stub */ },
                    shape = FloatingActionButtonDefaults.shape,
                    modifier = Modifier
                        .align(Alignment.Center)
                        .size(70.dp)
                        .offset(y = 50.dp)
                ) {
                    Icon(
                        painter = painterResource(R.drawable.fas_plus),
                        contentDescription = null,
                        modifier = Modifier.size(45.dp)
                    )
                }
            }
        },

        floatingActionButtonPosition = FabPosition.Center,

    ) { innerPadding ->
        Text(
            modifier = Modifier.padding(innerPadding),
            text = "Hello World"
        )
    }
}

@Preview(showBackground = true)
@Composable
fun GreetingPreview() {
    Simplecloudnotifier2Theme {

        Content()

    }
}