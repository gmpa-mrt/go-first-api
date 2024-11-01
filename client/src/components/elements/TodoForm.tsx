import { Button, Flex, Input, Spinner } from "@chakra-ui/react";
import {FormEvent, useState} from "react";
import { IoMdAdd } from "react-icons/io";
import {useMutation, useQueryClient} from "@tanstack/react-query"
import {BASE_URL} from "@/App.tsx";

const TodoForm = () => {
	const [newTodo, setNewTodo] = useState("");

	const queryClient = useQueryClient()

	const {mutate: createTodo, isPending: isCreating} = useMutation({
		mutationKey: ["todos"],
		mutationFn: async (e: FormEvent) => {
			e.preventDefault()
			try {
				const res = await fetch(`${BASE_URL}/todos`, {
					method: "POST",
					headers: {
						"Content-type" : "application/json "
					},
					body: JSON.stringify({body: newTodo})
				})

				const data = await res.json()

				setNewTodo("")

				if(!res.ok) {
					return new Error(data.error || "Something went wrong")
				}

				return data
			} catch (error) {
				console.log(error)
			}
		},
		onSuccess: () => queryClient.invalidateQueries({queryKey: ["todos"]}),
		onError: (error) => alert(error)
	})


	return (
		<form onSubmit={createTodo}>
			<Flex gap={2}>
				<Input
					type='text'
					value={newTodo}
					onChange={(e) => setNewTodo(e.target.value)}
					ref={(input) => input && input.focus()}
				/>
				<Button
					mx={2}
					type='submit'
					_active={{
						transform: "scale(.97)",
					}}
				>
					{isCreating ? <Spinner size={"xs"} /> : <IoMdAdd size={30} />}
				</Button>
			</Flex>
		</form>
	);
};
export default TodoForm;