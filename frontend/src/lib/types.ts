export interface LoginUser {
  email: string;
  password: string;
}

export interface RegisterUser extends LoginUser {
  name?: string;
  desiredJobType?: string;
  desiredLocation?: string;
  desiredCompanySize?: string;
  careerAxis1?: string;
  careerAxis2?: string;
  selfPr?: string;
}

export interface UserProfileUpdate {
  name: string;
  desiredJobType: string;
  desiredLocation: string;
  desiredCompanySize: string;
  careerAxis1: string;
  careerAxis2: string;
  selfPr: string;
}
