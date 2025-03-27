import { NextRequest, NextResponse } from 'next/server';
import jwt from 'jsonwebtoken';
import prisma from '@/lib/prisma';
import { RegisterUser } from '@/lib/types';

const JWT_SECRET = process.env.JWT_SECRET;

export async function POST(req: NextRequest) {
  const body: RegisterUser = await req.json();
  const {
    email,
    password,
    name,
    desiredJobType,
    desiredLocation,
    desiredCompanySize,
    careerAxis1,
    careerAxis2,
    selfPr,
  } = body;

  try {
    const user = await prisma.user.create({
      data: {
        email,
        password,
        name,
        desiredJobType,
        desiredLocation,
        desiredCompanySize,
        careerAxis1,
        careerAxis2,
        selfPr,
      },
    });
    if (JWT_SECRET) {
      const token = jwt.sign({ id: user.id, email: user.email }, JWT_SECRET, {
        expiresIn: '1h',
      });
      return NextResponse.json(
        { message: 'サインアップ成功', user, token },
        { status: 200 },
      );
    }
    return NextResponse.json(user);
  } catch (e) {
    return NextResponse.json(e);
  }
}
