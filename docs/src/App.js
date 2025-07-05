import React from 'react';
import { ForceGraph3D } from 'react-force-graph';
import { UnrealBloomPass } from 'three/examples/jsm/postprocessing/UnrealBloomPass';
import styled from 'styled-components';
import Button from 'react-bootstrap/Button';
import 'bootstrap/dist/css/bootstrap.min.css';
import backgroundData from './data.json';
import Stack from 'react-bootstrap/Stack';
import logoWhite from './logo-white.svg'
import { BsGithub } from "react-icons/bs";

import Nav from 'react-bootstrap/Nav';
import NavDropdown from 'react-bootstrap/NavDropdown';


// https://stackoverflow.com/a/38463360
const BackgroundGraph = styled.section`
width: 100vw;
height: 100vh;
display: block;
position: fixed;
top: 0;
left: 0;
z-index: -9999;
`;

function App() {

  const { useRef, useEffect } = React;

  const fgRef = useRef();

  useEffect(() => {
    const bloomPass = new UnrealBloomPass();
    bloomPass.strength = 1;
    bloomPass.radius = 1;
    bloomPass.threshold = 0;
    fgRef.current.postProcessingComposer().addPass(bloomPass);
  }, []);

  return (
    <>
      <style type="text/css">
        {`
          .dropdown>a::after {
            color: white;
          }
          .dropdown-menu {
            background-color: transparent;
            .dropdown-item:hover {
              text-decoration: underline;
              color: #fff;
              background-color: transparent;
            }
          }
        `}
      </style>

      {/* https://stackoverflow.com/a/69321054 */}
      <Stack direction="horizontal" gap={5} className='mt-5 mx-5'>
        <div className="p-2 text-white">
          <img src={logoWhite} alt='logo' />
          QubitPi
        </div>

        <Nav className="justify-content-center">
          <NavDropdown title={
            <span style={{ color: 'white' }}>HashiCorp</span>
          } id="basic-nav-dropdown">
            <NavDropdown.Item href="https://packer.qubitpi.org/packer">
              <span style={{ color: 'white' }}>Packer</span>
            </NavDropdown.Item>
            <NavDropdown.Item href="https://packer.qubitpi.org/terraform">
              <span style={{ color: 'white' }}>Terraform</span>
            </NavDropdown.Item>
          </NavDropdown>
        </Nav>

        <div className="p-2 ms-auto text-white">
          <a href="https://github.com/QubitPi/Henschel" style={{ color: '#ffffff' }}>
            <BsGithub
              size={30}
              onMouseOver={({ target }) => target.style.color = '#0C6EFD'}
              onMouseOut={({ target }) => target.style.color = '#ffffff'}
            />
          </a>
        </div>
        <Button variant="outline-light" className="p-2" href="https://github.com/QubitPi/Henschel">Start Free</Button>
      </Stack>

      {/* https://stackoverflow.com/a/52284399 */}
      <Button variant="outline-primary" size="lg" href="https://github.com/QubitPi/Henschel" style={{
        position: 'absolute', left: '50%', top: '70%',
        transform: 'translate(-50%, -50%)'
      }}>Henschel</Button>

      <BackgroundGraph>
        <ForceGraph3D
          ref={fgRef}
          backgroundColor="#000003"
          graphData={backgroundData}
          nodeLabel={node => `${node.user}: ${node.description}`}
          nodeAutoColorBy="user"
          linkDirectionalParticles={1}
        />
      </BackgroundGraph>
    </>
  );
}

export default App;
