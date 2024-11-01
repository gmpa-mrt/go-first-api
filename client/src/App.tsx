import {Container, Stack} from "@chakra-ui/react";
import Navbar from "@/components/elements/Navbar.tsx";
import TodoForm from "@/components/elements/TodoForm.tsx";
import TodoList from "@/components/elements/TodoList.tsx";

export const BASE_URL = import.meta.env.MODE === "develop" ? "http://localhost:4000/api" : "/api"

function App() {
  return (
    <Stack h={'100vh'}>
      <Navbar />
      <Container>
          <TodoForm />
          <TodoList />
      </Container>
    </Stack>
  )
}

export default App
