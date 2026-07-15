import React, { useState, useEffect } from 'react';
import {
  ThemeProvider,
  createTheme,
  CssBaseline,
  Container,
  Box,
  Typography,
  Button,
  Card,
  CardContent,
  Stack,
  IconButton,
  AppBar,
  Toolbar,
  Paper
} from '@mui/material';
import {
  LocalDrink as JuiceIcon,
  TouchApp as TouchIcon,
  GitHub as GitHubIcon,
  LightMode,
  DarkMode,
  Code as CodeIcon
} from '@mui/icons-material';
import { SweetJuice } from './bridge';

const theme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#f2ca50',
    },
    background: {
      default: '#111415',
      paper: '#1a1d1e',
    },
  },
  typography: {
    fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
    h4: {
      fontWeight: 800,
    },
  },
  shape: {
    borderRadius: 16,
  },
});

function App() {
  const [count, setCount] = useState(0);
  const [status, setStatus] = useState('Waiting for bridge...');
  const [isJuiced, setIsJuiced] = useState(false);

  useEffect(() => {
    // Test bridge connection
    const testBridge = async () => {
      try {
        const result = await SweetJuice.CallGo('Ping', 'Frontend ready');
        if (result.error) {
          setStatus('Ready (Browser Mode)');
        } else {
          setStatus('Connected to Go Backend');
        }
      } catch (e) {
        setStatus('Offline');
      }
    };

    testBridge();

    // Listen for events
    SweetJuice.on('backend_event', (data) => {
      console.log('Received event from Go:', data);
    });
  }, []);

  const handleJuice = () => {
    setIsJuiced(true);
    setCount(c => c + 1);
    SweetJuice.CallGo('Action', 'Juiced!');
    setTimeout(() => setIsJuiced(false), 500);
  };

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Box sx={{ flexGrow: 1, minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
        <AppBar position="static" color="transparent" elevation={0} sx={{ borderBottom: '1px solid rgba(255,255,255,0.05)' }}>
          <Toolbar>
            <JuiceIcon sx={{ mr: 2, color: 'primary.main' }} />
            <Typography variant="h6" component="div" sx={{ flexGrow: 1, fontWeight: 700, letterSpacing: -0.5 }}>
              SWEET JUICE
            </Typography>
            <IconButton color="inherit" onClick={() => window.open('https://github.com/sweet-juice/sweetjuice', '_blank')}>
              <GitHubIcon />
            </IconButton>
          </Toolbar>
        </AppBar>

        <Container maxWidth="sm" sx={{ mt: 8, mb: 4, flexGrow: 1 }}>
          <Stack spacing={4} alignItems="center">

            <Box sx={{ position: 'relative' }}>
                <Paper
                    elevation={24}
                    sx={{
                        width: 120,
                        height: 120,
                        borderRadius: '50%',
                        display: 'flex',
                        justifyContent: 'center',
                        alignItems: 'center',
                        background: 'linear-gradient(45deg, #f2ca50 30%, #ff8e53 90%)',
                        transition: 'transform 0.2s',
                        transform: isJuiced ? 'scale(1.2)' : 'scale(1)',
                    }}
                >
                    <JuiceIcon sx={{ fontSize: 60, color: '#111415' }} />
                </Paper>
            </Box>

            <Box sx={{ textAlign: 'center' }}>
              <Typography variant="h4" gutterBottom>
                Modern Go Mobile
              </Typography>
              <Typography variant="body1" color="text.secondary" sx={{ maxWidth: 300, mx: 'auto' }}>
                Build high-performance native apps with a Go backend and React frontend.
              </Typography>
            </Box>

            <Card sx={{ width: '100%', bgcolor: 'background.paper', border: '1px solid rgba(255,255,255,0.05)' }}>
              <CardContent>
                <Stack spacing={2}>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Typography variant="overline" color="text.secondary">
                      Bridge
                    </Typography>
                    <Box sx={{ px: 1.5, py: 0.5, borderRadius: 10, bgcolor: 'rgba(242, 202, 80, 0.1)', color: 'primary.main', fontSize: '0.75rem', fontWeight: 700 }}>
                      {status}
                    </Box>
                  </Box>

                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Typography variant="h6">
                      Press Count
                    </Typography>
                    <Typography variant="h4" color="primary">
                      {count}
                    </Typography>
                  </Box>

                  <Button
                    fullWidth
                    variant="contained"
                    size="large"
                    startIcon={<TouchIcon />}
                    onClick={handleJuice}
                    sx={{
                        height: 56,
                        fontWeight: 700,
                        boxShadow: '0 8px 16px rgba(242, 202, 80, 0.2)'
                    }}
                  >
                    Tap the Juice
                  </Button>
                </Stack>
              </CardContent>
            </Card>

            <Stack direction="row" spacing={1}>
                <CodeIcon fontSize="small" color="disabled" />
                <Typography variant="caption" color="text.disabled">
                    AppTemplate/frontend/src/App.jsx
                </Typography>
            </Stack>

          </Stack>
        </Container>
      </Box>
    </ThemeProvider>
  );
}

export default App;
