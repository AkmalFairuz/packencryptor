import {
    Box,
    Button,
    ChakraProvider,
    extendTheme,
    Input,
    InputGroup,
    InputRightElement,
    useToast
} from '@chakra-ui/react'
import {useState} from "react";
import {EncryptPack, SelectDestDialog, SelectPathDialog} from "../wailsjs/go/main/App.js";

const theme = extendTheme({
    styles: {
        global: (props) => ({
            body: {
                bg: "#181818",
                color: "white",
            }
        })
    }
})

function App() {
    const toast = useToast()
    const [path, setPath] = useState("")
    const [dest, setDest] = useState("")
    const [key, setKey] = useState("")
    const [loading, setLoading] = useState(false)
    const handlePathDialog = () => {
        SelectPathDialog().then((r) => {
            if (r === "") {
                return
            }
            setPath(r)
        })
    }
    const handleDestDialog = () => {
        SelectDestDialog().then((r) => {
            if (r === "") {
                return
            }
            setDest(r)
        })
    }
    const handleEncrypt = () => {
        if (dest === "" || path === "") {
            toast({
                position: "bottom-right",
                title: "Failed",
                status: "error",
                description: "Please select pack and destination",
                duration: 5000,
                isClosable: true,
            })
            return
        }
        setLoading(true)
        EncryptPack(path, dest).then((key) => {
            toast({
                position: "bottom-right",
                title: "Pack Encrypted",
                status: "success",
                duration: 5000,
                isClosable: true,
            })
            setKey(key)
        }).catch((e) => {
            toast({
                position: "bottom-right",
                description: e,
                status: "error",
                duration: 10000,
                isClosable: true,
            })
        }).finally(() => {
            setLoading(false)
        })
    }
    const copyKey = () => {
        if (key === "") {
            return
        }
        navigator.clipboard.writeText(key).then(() => {
            toast({
                position: "bottom-right",
                title: "Key copied!",
                status: "success",
                duration: 2000,
                isClosable: true,
            })
        })
    }
    return <ChakraProvider theme={theme}>
        <Box p={4} pt={1}>
            <Box textAlign={"center"} fontSize={"4xl"} fontWeight={"bold"} mb={2}>
                PackEncryptor
            </Box>
            <Box mx={2} mb={3}>
                <Box mb={2}>
                    <Box as={"span"} fontWeight={"semibold"}>Pack:</Box> {path}
                </Box>
                <Button onClick={handlePathDialog} size={"sm"} colorScheme={"blue"}>Select pack</Button>
                <Box mt={1} fontSize={"sm"} textColor={"gray"}>
                    A Minecraft Bedrock resource pack file path (.zip or .mcpack)
                </Box>
            </Box>
            <Box mx={2} mb={3}>
                <Box mb={2}>
                    <Box as={"span"} fontWeight={"semibold"}>Destination:</Box> {dest}
                </Box>
                <Button onClick={handleDestDialog} size={"sm"} colorScheme={"blue"}>Select destination</Button>
                <Box mt={1} fontSize={"sm"} textColor={"gray"}>
                    The destination where the encrypted pack will be saved
                </Box>
            </Box>
            <Box mx={2} mb={3}>
                <Box mb={2}>
                    <Box as={"span"} fontWeight={"semibold"}>Key:</Box>
                </Box>
                <InputGroup size={"md"}>
                    <Input
                        fontFamily={"monospace"}
                        pr="4.5rem"
                        value={key}
                    />
                    <InputRightElement width={"4.5rem"}>
                        <Button h={"1.75rem"} size={"sm"} colorScheme={"blue"} onClick={copyKey}>
                            Copy
                        </Button>
                    </InputRightElement>
                </InputGroup>
                <Box mt={1} mb={2} fontSize={"sm"} textColor={"gray"}>
                    The encryption key to be used in your server configuration
                </Box>
            </Box>
            <Box mx={2} mb={3}>
                <Button w={"full"} colorScheme={"green"} onClick={handleEncrypt} isLoading={loading}>Encrypt</Button>
            </Box>
        </Box>
    </ChakraProvider>
}

export default App
