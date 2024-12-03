// Copyright 2024, Horizen Labs, Inc.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

const { utils, getCurveFromName, Scalar } = require("ffjavascript");
const { unstringifyBigInts } = utils;
const { toRprLE } = Scalar;

async function convertProof(proofJSON, gnarkCurveName) {
  const proof = unstringifyBigInts(proofJSON);
  const curve = await getCurveFromName(gnarkCurveName);
  const curveName = getCurveName(curve);
  let endianess;
  switch (curveName) {
    case "Bn254":
      endianess = "LE";
      break;
    case "Bls12_381":
      endianess = "BE";
      break;
  }
  return {
    a: pointToHex(curve.G1, [proof.Ar.X, proof.Ar.Y], endianess),
    b: pointToHex(curve.G2, [[proof.Bs.X.A0, proof.Bs.X.A1], [proof.Bs.Y.A0, proof.Bs.Y.A1]], endianess),
    c: pointToHex(curve.G1, [proof.Krs.X, proof.Krs.Y], endianess),
  };
}

async function convertVk(vkJSON, gnarkCurveName) {
  const vk = unstringifyBigInts(vkJSON);
  const curve = await getCurveFromName(gnarkCurveName);
  const curveName = getCurveName(curve);
  let endianess;
  switch (curveName) {
    case "Bn254":
      endianess = "LE";
      break;
    case "Bls12_381":
      endianess = "BE";
      break;
  }
  return {
    curve: curveName,
    alphaG1: pointToHex(curve.G1, [vk.G1.Alpha.X, vk.G1.Alpha.Y], endianess),
    betaG2: pointToHex(curve.G2, [[vk.G2.Beta.X.A0, vk.G2.Beta.X.A1], [vk.G2.Beta.Y.A0, vk.G2.Beta.Y.A1]], endianess),
    gammaG2: pointToHex(curve.G2, [[vk.G2.Gamma.X.A0, vk.G2.Gamma.X.A1], [vk.G2.Gamma.Y.A0, vk.G2.Gamma.Y.A1]], endianess),
    deltaG2: pointToHex(curve.G2, [[vk.G2.Delta.X.A0, vk.G2.Delta.X.A1], [vk.G2.Delta.Y.A0, vk.G2.Delta.Y.A1]], endianess),
    gammaAbcG1: vk.G1.K.map((point) => pointToHex(curve.G1, [point.X, point.Y], endianess)),
  };
}

async function convertPub(pubJSON, gnarkCurveName) {
  const pub = unstringifyBigInts(pubJSON);
  const curve = await getCurveFromName(gnarkCurveName);
  return pub.map((scalar) => scalarToHexLE(curve.Fr, scalar));
}

function getCurveName(curve) {
  let curveName;
  switch (curve.name) {
    case "bn128":
      curveName = "Bn254";
      break;
    case "bls12381":
      curveName = "Bls12_381";
      break;
    default:
      throw new Error(`Curve ${curve.name} is not supported`);
  }
  return curveName;
}

const pointToHex = (curveType, point, endianess) => {
  const p = curveType.fromObject(point);
  const pUncompressed = curveType.toUncompressed(p);
  let x = pUncompressed.slice(0, curveType.F.n8);
  let y = pUncompressed.slice(curveType.F.n8);
  switch (endianess) {
    case "LE":
      x.reverse();
      y.reverse();
      break;
    case "BE":
      break;
    default:
      throw new Error("pointEndianess must be either BE or LE");
  }
  return "0x" + toHexString(x) + toHexString(y);
};

function scalarToHexLE(scalarType, scalar) {
  const buff = new Uint8Array(scalarType.n8);
  toRprLE(buff, 0, scalar, buff.byteLength);
  return "0x" + toHexString(buff);
}

function toHexString(byteArray) {
  return Array.from(byteArray, function (byte) {
    return ("0" + (byte & 0xff).toString(16)).slice(-2);
  }).join("");
}

module.exports = {
  convertProof,
  convertVk,
  convertPub,
  pointToHexLE: pointToHex,
};
