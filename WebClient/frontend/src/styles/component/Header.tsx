import styled from "styled-components";
import { theme } from "../theme";

export const Header = styled.header`
  position: fixed;
  top: 0;
  left: 0;
  height: 56px;
  width: 100%;
  background-color: ${theme.colors.background};
  border-bottom: 1px solid ${theme.colors.border};
  z-index: 100;

  display: flex;
  align-items: center;
  justify-content: center;
`;
