import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { MantineProvider } from "@mantine/core";
import "@mantine/core/styles.css";

import Uploady from "@rpldy/uploady";

const queryClient = new QueryClient();
const root = ReactDOM.createRoot(document.getElementById('root'));

root.render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <MantineProvider>
       <Uploady
    destination={{
        url: "https://хз/пока"
    }}
    autoUpload={false}
    accept=".pdf"
>
          <App />
        </Uploady>
      </MantineProvider>
    </QueryClientProvider>
  </React.StrictMode>
);