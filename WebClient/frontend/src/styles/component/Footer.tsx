import styled from "styled-components";
import { theme } from "../theme";

export const Footer = styled.footer`
  position: fixed;
  bottom: 0;
  left: 0;
  height: 32px;
  width: 100%;
  background-color: ${theme.colors.background};
  border-top: 1px solid ${theme.colors.border};

  display: flex;
  align-items: center;
  justify-content: center;
`;
