import styled from "styled-components";
import { theme } from "../theme";

export const Overlay = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
`;

export const Modal = styled.div`
  background: ${theme.colors.surface};
  padding: 20px;
  border-radius: ${theme.radius.md};
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-width: 300px;
  box-shadow: ${theme.shadows.card};
`;
