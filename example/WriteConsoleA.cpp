#include <iostream>

#ifdef _WIN32
#define WIN32_LEAN_AND_MEAN
#include <Windows.h>
#endif

int main()
{
#ifdef _WIN32
	std::string value = "this is an example output using WriteConsoleA that cannot be piped or redirected";

	HANDLE hConsole = GetStdHandle(STD_OUTPUT_HANDLE);
	WriteConsoleA(hConsole, value.data(), value.size(), NULL, NULL);

	return 0;
#else
	std::cout << "this test only works on windows" << std::endl;
	return 1;
#endif
}

