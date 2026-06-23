import styled, { keyframes } from "styled-components";
import { Center } from "../styles/layout/Center";
import { Stack } from "../styles/primitives/Stack";

const spin = keyframes`
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
`;

const Spinner = styled.div`
  width: 40px;
  height: 40px;
  border: 4px solid #ccc;
  border-top: 4px solid #333;
  border-radius: 50%;
  animation: ${spin} 1s linear infinite;
`;

const Loader = () => {
  return (
    <Center>
      <Stack>
        <Spinner />
      </Stack>
    </Center>
  );
};

export default Loader;
