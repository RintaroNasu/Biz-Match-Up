import { NextRequest, NextResponse } from 'next/server';
import prisma from '../../../../lib/prisma';

export async function POST(
  req: NextRequest,
  { params }: { params: { userId: string } },
) {
  const userId = (await params).userId;
  const parsedUserId = parseInt(userId);
  const body = await req.json();

  try {
    const updatedUser = await prisma.user.update({
      where: { id: parsedUserId },
      data: body,
    });

    return NextResponse.json({ success: true, user: updatedUser });
  } catch (error) {
    return NextResponse.json({ success: false, error });
  }
}
